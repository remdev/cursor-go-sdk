# Vendored dependencies

## `@anysphere/proto`

Connect/protobuf message definitions for the SDK bridge wire protocol. Not published on the public npm registry; vendored here so `@cursor-go-sdk/cursor-sdk-bridge` installs self-contained via `npm install`.

Source: same layout as the bridge bundled in the official `cursor-sdk` Python wheel.

When regenerating, copy from an installed `cursor_sdk._vendor.bridge` tree or update from upstream proto exports if Cursor publishes them.
