# Publishing `@cursor-go-sdk/cursor-sdk-bridge`

## One-time setup

1. npm account is a member of the **`cursor-go-sdk`** org with publish rights.
2. Enable **2FA** on npm (recommended; may be required for publish).
3. Log in locally:

```bash
npm login
npm whoami
```

## Pre-publish checklist

```bash
cd bridge
npm ci
npm run build
./bin/cursor-sdk-bridge --help
npm publish --dry-run --access public
```

Review the tarball contents: `bin/`, `dist/`, `gen/`, `proto/` must be included.

Regenerate after `.proto` changes:

```bash
npm run generate && npm run build
```

## Publish

```bash
cd bridge
npm publish --access public
```

`publishConfig.access` in `package.json` is already `"public"`.

## After publish

```bash
npm install -g @cursor-go-sdk/cursor-sdk-bridge
cursor-sdk-bridge --help
```

Package page: https://www.npmjs.com/package/@cursor-go-sdk/cursor-sdk-bridge

## Version bumps

1. Edit `"version"` in `bridge/package.json` (semver).
2. `npm install` to refresh `package-lock.json` if dependencies changed.
3. Commit, tag optional (`cursor-sdk-bridge-v1.0.1`), publish.

## CI token (later)

Create an npm **Granular Access Token** scoped to `@cursor-go-sdk/*` with publish permission. Store as `NPM_TOKEN` in GitHub Actions.

```yaml
- uses: actions/setup-node@v4
  with:
    node-version: "20"
    registry-url: https://registry.npmjs.org
- run: npm publish --access public
  working-directory: bridge
  env:
    NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
```
