//go:build embedbridge

package bridge

import _ "embed"

// prebuilt/bridge.tar.gz is produced by cmd/fetch-bridge and is not committed.
// Build with -tags embedbridge after running fetch-bridge to ship the bridge
// inside the module binary (still extracted to the cache directory on first use).
//
//go:embed prebuilt/bridge.tar.gz
var embeddedBridgeTarballData []byte

func embeddedBridgeTarball() []byte {
	return embeddedBridgeTarballData
}
