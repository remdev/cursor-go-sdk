package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remdev/cursor-go-sdk/examples/agentutil"
)

type lineKind int

const (
	lineUser lineKind = iota
	lineAssistant
	lineMeta
	lineError
)

type transcriptLine struct {
	kind lineKind
	text string
}

type uiMode int

const (
	modeInput uiMode = iota
	modeModelPicker
)

var (
	headerStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	metaStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	userStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("117"))
	assistantStyle = lipgloss.NewStyle()
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	inputStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("241"))
	pickerStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	pickerActive   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	pickerItem     = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	busyStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
)

type appModel struct {
	sess       *session
	ctx        context.Context
	cancel     context.CancelFunc
	viewport   viewport.Model
	input      textinput.Model
	lines      []transcriptLine
	busy       bool
	mode       uiMode
	models     []modelChoice
	modelIndex int
	width      int
	height     int
	quitting   bool
	streaming  bool
	runEvents  <-chan agentEvent
	runDone    <-chan error
}

func newAppModel(sess *session) *appModel {
	ti := textinput.New()
	ti.Placeholder = "Ask a question or type /help"
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 80

	m := &appModel{
		sess:     sess,
		viewport: viewport.New(80, 20),
		input:    ti,
	}
	m.appendMeta(welcomeText(sess))
	m.refreshViewport()
	return m
}

func welcomeText(sess *session) string {
	return fmt.Sprintf(
		"Lightweight coding agent · %s · %s\nType a prompt or /help for commands.",
		sess.mode,
		sess.executionTarget(),
	)
}

func (m *appModel) Init() tea.Cmd {
	return textinput.Blink
}

type runEventMsg agentEvent

type runFinishedMsg struct {
	err error
}

type runEventsClosedMsg struct{}

type modelsLoadedMsg struct {
	models []modelChoice
	err    error
}

func (m *appModel) ensureCtx() context.Context {
	if m.ctx == nil {
		m.ctx, m.cancel = context.WithCancel(context.Background())
	}
	return m.ctx
}

func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout()
		return m, nil

	case runEventsClosedMsg:
		return m, nil

	case runEventMsg:
		m.applyEvent(agentEvent(msg))
		m.refreshViewport()
		if m.runEvents != nil {
			return m, waitRunEvent(m.runEvents)
		}
		return m, nil

	case runFinishedMsg:
		m.busy = false
		m.streaming = false
		m.runEvents = nil
		m.runDone = nil
		if msg.err != nil {
			m.appendError(msg.err.Error())
			m.refreshViewport()
		}
		return m, nil

	case modelsLoadedMsg:
		if msg.err != nil {
			m.appendError(msg.err.Error())
			m.refreshViewport()
			return m, nil
		}
		m.models = msg.models
		m.modelIndex = 0
		for i, c := range m.models {
			if modelSelectionKey(c.value) == modelSelectionKey(m.sess.model) {
				m.modelIndex = i
				break
			}
		}
		m.mode = modeModelPicker
		m.input.Blur()
		return m, nil

	case tea.KeyMsg:
		if m.quitting {
			return m, tea.Quit
		}
		switch m.mode {
		case modeModelPicker:
			return m.updateModelPicker(msg)
		default:
			return m.updateInput(msg)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func waitRunEvent(ch <-chan agentEvent) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-ch
		if !ok {
			return runEventsClosedMsg{}
		}
		return runEventMsg(ev)
	}
}

func waitRunDone(ch <-chan error) tea.Cmd {
	return func() tea.Msg {
		err, ok := <-ch
		if !ok {
			return runFinishedMsg{}
		}
		return runFinishedMsg{err: err}
	}
}

func (m *appModel) updateModelPicker(msg tea.KeyMsg) (*appModel, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.mode = modeInput
		m.input.Focus()
		return m, textinput.Blink
	case "up", "k":
		if m.modelIndex > 0 {
			m.modelIndex--
		}
	case "down", "j":
		if m.modelIndex < len(m.models)-1 {
			m.modelIndex++
		}
	case "enter":
		if len(m.models) > 0 {
			choice := m.models[m.modelIndex]
			m.sess.setModel(choice.value)
			m.appendMeta(fmt.Sprintf("Model set to %s.", choice.label))
			m.refreshViewport()
		}
		m.mode = modeInput
		m.input.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

func (m *appModel) updateInput(msg tea.KeyMsg) (*appModel, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.busy {
			if err := m.sess.cancelCurrentRun(m.ensureCtx()); err != nil {
				switch {
				case errors.Is(err, ErrNoActiveRun):
					m.appendMeta("No active run to cancel.")
				case errors.Is(err, ErrRunNotCancelable):
					m.appendMeta("This run cannot be cancelled.")
				default:
					m.appendError(err.Error())
				}
			} else {
				m.appendMeta("Run cancelled.")
			}
			m.refreshViewport()
			return m, nil
		}
		m.quitting = true
		if m.cancel != nil {
			m.cancel()
		}
		return m, tea.Quit

	case "pgup":
		m.viewport.LineUp(3)
		return m, nil
	case "pgdown":
		m.viewport.LineDown(3)
		return m, nil

	case "enter":
		line := strings.TrimSpace(m.input.Value())
		if line == "" || m.busy {
			return m, nil
		}
		m.input.SetValue("")
		if cmd := m.handleSubmit(line); cmd != nil {
			return m, cmd
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *appModel) handleSubmit(line string) tea.Cmd {
	if cmdName := getSlashCommand(line); cmdName != "" {
		return m.handleSlash(cmdName)
	}
	m.appendUser(line)
	m.refreshViewport()
	return m.startRun(line)
}

func (m *appModel) handleSlash(cmd string) tea.Cmd {
	switch cmd {
	case "/help", "/?":
		m.appendMeta(formatSlashHelp())
	case "/exit", "/quit":
		m.quitting = true
		if m.cancel != nil {
			m.cancel()
		}
		return tea.Quit
	case "/reset":
		if m.busy {
			m.appendMeta("Wait for the current run to finish before resetting.")
		} else if err := m.sess.reset(m.ensureCtx()); err != nil {
			m.appendError(err.Error())
		} else {
			m.lines = nil
			m.appendMeta("Agent reset. Context cleared.")
			m.appendMeta(welcomeText(m.sess))
		}
	case "/local":
		if m.busy {
			m.appendMeta("Wait for the current run to finish before switching mode.")
		} else if err := m.sess.setExecutionMode(m.ensureCtx(), modeLocal); err != nil {
			m.appendError(err.Error())
		} else {
			m.appendMeta("Local mode enabled. Target: " + m.sess.executionTarget())
		}
	case "/cloud":
		if m.busy {
			m.appendMeta("Wait for the current run to finish before switching mode.")
		} else if err := m.sess.setExecutionMode(m.ensureCtx(), modeCloud); err != nil {
			m.appendError(err.Error())
		} else {
			m.appendMeta("Cloud mode enabled. Target: " + m.sess.executionTarget())
		}
	case "/model":
		if m.busy {
			m.appendMeta("Wait for the current run to finish before changing model.")
		} else {
			m.refreshViewport()
			return m.loadModels()
		}
	default:
		m.appendMeta("Unknown command. Type /help.")
	}
	m.refreshViewport()
	return nil
}

func (m *appModel) loadModels() tea.Cmd {
	ctx := m.ensureCtx()
	sess := m.sess
	return func() tea.Msg {
		models, err := sess.listModels(ctx)
		return modelsLoadedMsg{models: models, err: err}
	}
}

// startRun sets run state synchronously in Update before cmds are scheduled.
func (m *appModel) startRun(prompt string) tea.Cmd {
	ctx := m.ensureCtx()
	sess := m.sess
	events := make(chan agentEvent, 64)
	done := make(chan error, 1)

	m.busy = true
	m.streaming = false
	m.runEvents = events
	m.runDone = done

	go func() {
		err := sess.sendPrompt(ctx, prompt, func(ev agentEvent) {
			events <- ev
		})
		close(events)
		done <- err
	}()

	return tea.Batch(waitRunEvent(events), waitRunDone(done))
}

func (m *appModel) applyEvent(ev agentEvent) {
	switch ev.kind {
	case eventAssistantDelta:
		if ev.text == "" {
			return
		}
		if m.streaming && len(m.lines) > 0 && m.lines[len(m.lines)-1].kind == lineAssistant {
			m.lines[len(m.lines)-1].text += ev.text
		} else {
			m.lines = append(m.lines, transcriptLine{kind: lineAssistant, text: ev.text})
			m.streaming = true
		}
	case eventThinking:
		m.streaming = false
		m.appendMeta("[thinking] " + ev.text)
	case eventTool:
		m.streaming = false
		line := fmt.Sprintf("[tool] %s %s", ev.toolStatus, ev.toolName)
		if ev.toolParams != "" {
			line += " " + ev.toolParams
		}
		m.appendMeta(line)
	case eventStatus:
		m.streaming = false
		line := "[status] " + ev.status
		if ev.statusMessage != "" {
			line += " " + ev.statusMessage
		}
		m.appendMeta(line)
	case eventTask:
		m.streaming = false
		m.appendMeta("[task] " + agentutil.Compact(ev.taskStatus+" "+ev.taskText))
	case eventResult:
		m.streaming = false
		m.appendMeta(fmt.Sprintf("[done] status=%s durationMs=%d", ev.resultStatus, ev.durationMS))
	case eventError:
		m.streaming = false
		if ev.err != nil {
			m.appendError(ev.err.Error())
		}
	}
}

func (m *appModel) appendUser(text string) {
	m.streaming = false
	m.lines = append(m.lines, transcriptLine{kind: lineUser, text: text})
}

func (m *appModel) appendMeta(text string) {
	m.streaming = false
	m.lines = append(m.lines, transcriptLine{kind: lineMeta, text: text})
}

func (m *appModel) appendError(text string) {
	m.streaming = false
	m.lines = append(m.lines, transcriptLine{kind: lineError, text: text})
}

func (m *appModel) refreshViewport() {
	var b strings.Builder
	for _, line := range m.lines {
		switch line.kind {
		case lineUser:
			b.WriteString(userStyle.Render("> " + line.text))
		case lineAssistant:
			b.WriteString(assistantStyle.Render(line.text))
		case lineMeta:
			b.WriteString(metaStyle.Render(line.text))
		case lineError:
			b.WriteString(errorStyle.Render(line.text))
		}
		b.WriteByte('\n')
	}
	m.viewport.SetContent(b.String())
	m.viewport.GotoBottom()
}

func (m *appModel) layout() {
	headerH := 2
	inputH := 3
	pickerH := 0
	if m.mode == modeModelPicker {
		pickerH = min(8, len(m.models)+3)
		if pickerH < 5 {
			pickerH = 5
		}
	}
	bodyH := m.height - headerH - inputH - pickerH
	if bodyH < 3 {
		bodyH = 3
	}
	m.viewport.Width = m.width
	m.viewport.Height = bodyH
	m.input.Width = m.width - 4
}

func (m *appModel) View() string {
	if m.quitting {
		return ""
	}

	header := headerStyle.Render(fmt.Sprintf(
		"%s · model %s · target %s",
		strings.ToUpper(string(m.sess.mode)),
		m.sess.modelLabel(),
		m.sess.executionTarget(),
	))
	if m.busy {
		header += busyStyle.Render(" · running…")
	}

	var b strings.Builder
	b.WriteString(header)
	b.WriteByte('\n')
	b.WriteString(m.viewport.View())
	b.WriteByte('\n')

	if m.mode == modeModelPicker {
		b.WriteString(m.renderModelPicker())
		b.WriteByte('\n')
	}

	prompt := "> "
	if m.busy {
		prompt = "… "
	}
	b.WriteString(inputStyle.Render(prompt + m.input.View()))
	return b.String()
}

func (m *appModel) renderModelPicker() string {
	var b strings.Builder
	b.WriteString(pickerStyle.Render("Select model (↑/↓, Enter, Esc)"))
	b.WriteByte('\n')
	for i, choice := range m.models {
		prefix := "  "
		style := pickerItem
		if i == m.modelIndex {
			prefix = "▸ "
			style = pickerActive
		}
		b.WriteString(style.Render(prefix + choice.label))
		b.WriteByte('\n')
	}
	return b.String()
}
