package bridge_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/remdev/cursor-go-sdk/internal/bridge"
)

func TestResolvePathFromNpmLayout020(t *testing.T) {
	root := t.TempDir()
	pkg := filepath.Join(root, "lib", "node_modules", "@cursor-go-sdk", "cursor-sdk-bridge")
	distDir := filepath.Join(pkg, "dist", "bin")
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		t.Fatal(err)
	}
	js := filepath.Join(distDir, "cursor-sdk-bridge.js")
	if err := os.WriteFile(js, []byte("#!/usr/bin/env node\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkg, "package.json"), []byte(`{"name":"@cursor-go-sdk/cursor-sdk-bridge","version":"0.0.2"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	globalBin := filepath.Join(root, "bin", "cursor-sdk-bridge")
	if err := os.MkdirAll(filepath.Dir(globalBin), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(js, globalBin); err != nil {
		t.Fatal(err)
	}

	t.Setenv("PATH", filepath.Join(root, "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("CURSOR_SDK_BRIDGE_BIN", "")
	t.Setenv("CURSOR_SDK_BRIDGE_ROOT", "")

	resolved, err := bridge.ResolvePath()
	if err != nil {
		t.Fatal(err)
	}
	got, _ := filepath.EvalSymlinks(resolved)
	want, _ := filepath.EvalSymlinks(js)
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestResolvePathRejectsLegacyShellLauncher(t *testing.T) {
	root := t.TempDir()
	pkg := filepath.Join(root, "lib", "node_modules", "@cursor-go-sdk", "cursor-sdk-bridge")
	binDir := filepath.Join(pkg, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	shell := filepath.Join(binDir, "cursor-sdk-bridge")
	if err := os.WriteFile(shell, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkg, "package.json"), []byte(`{"name":"@cursor-go-sdk/cursor-sdk-bridge","version":"0.0.1"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	globalBin := filepath.Join(root, "bin", "cursor-sdk-bridge")
	if err := os.MkdirAll(filepath.Dir(globalBin), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(shell, globalBin); err != nil {
		t.Fatal(err)
	}

	t.Setenv("PATH", filepath.Join(root, "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("CURSOR_SDK_BRIDGE_BIN", "")
	t.Setenv("CURSOR_SDK_BRIDGE_ROOT", "")

	_, err := bridge.ResolvePath()
	if err == nil {
		t.Fatal("expected error for legacy 0.0.1 shell launcher")
	}
	if !strings.Contains(err.Error(), "0.0.1") || !strings.Contains(err.Error(), "0.0.2") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveLaunchArgvUsesNodeForJS(t *testing.T) {
	repoBridge, err := filepath.Abs(filepath.Join("..", "..", "bridge"))
	if err != nil {
		t.Fatal(err)
	}
	js := filepath.Join(repoBridge, "dist", "bin", "cursor-sdk-bridge.js")
	if _, err := os.Stat(js); err != nil {
		t.Skip("bridge dist not built")
	}
	t.Setenv("CURSOR_SDK_BRIDGE_BIN", js)
	t.Setenv("CURSOR_SDK_BRIDGE_ROOT", "")
	t.Setenv("CURSOR_SDK_NODE_BIN", "/usr/bin/node")

	argv, err := bridge.ResolveLaunchArgvForTest()
	if err != nil {
		t.Fatal(err)
	}
	if len(argv) != 2 {
		t.Fatalf("argv: %v", argv)
	}
	if argv[0] != "/usr/bin/node" {
		t.Fatalf("expected node, got %q", argv[0])
	}
	if argv[1] != js {
		t.Fatalf("got js %q", argv[1])
	}
}

func TestResolvePathFromGlobalBinOnPATH(t *testing.T) {
	t.Setenv("CURSOR_SDK_BRIDGE_BIN", "")
	t.Setenv("CURSOR_SDK_BRIDGE_ROOT", "")

	_, err := bridge.ResolvePath()
	if err != nil {
		t.Skip(err)
	}
}
