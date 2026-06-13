# Bridge wire protocol (protobuf)

Connect RPC schema for **cursor-go-sdk** ↔ **cursor-sdk-bridge**.

- **Package:** `sdk.v1` (wire paths like `sdk.v1.SdkAgentService/CreateAgent`)
- **Semantics:** mapped to [`@cursor/sdk`](https://www.npmjs.com/package/@cursor/sdk) in `bridge/src/`
- **Go client:** hand-written JSON in `internal/connect/` (no Go codegen yet)

## Generate TypeScript

```bash
cd bridge
npm install
npm run generate && npm run build
```

Outputs to `bridge/gen/ts/` via [Buf](https://buf.build) + `@bufbuild/protoc-gen-es` + `@connectrpc/protoc-gen-connect-es`.

## Layout

| File | Service / types |
|------|-----------------|
| `sdk_messages.proto` | Shared messages (AgentOptions, RunStream, …) |
| `sdk_agent_service.proto` | Agent + run RPCs |
| `sdk_cursor_service.proto` | Me, ListModels, ListRepositories |
| `sdk_bridge_control_service.proto` | Ping, Shutdown, GetVersion, SetToolCallback |
| `sdk_custom_tool_callback_service.proto` | Go custom-tool callback client |
| `sdk_errors.proto` | Connect error details |

Owned and maintained in this repository.
