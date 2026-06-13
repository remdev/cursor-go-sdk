package bridge

import (
	"fmt"
	"runtime"
)

type platformSpec struct {
	ManifestOS         string
	ManifestArch       string
	NpmPlatformPackage string
}

func currentPlatform() (platformSpec, error) {
	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return platformSpec{
				ManifestOS:         "darwin",
				ManifestArch:       "arm64",
				NpmPlatformPackage: "@cursor/sdk-darwin-arm64",
			}, nil
		case "amd64":
			return platformSpec{
				ManifestOS:         "darwin",
				ManifestArch:       "amd64",
				NpmPlatformPackage: "@cursor/sdk-darwin-x64",
			}, nil
		}
	case "linux":
		switch runtime.GOARCH {
		case "arm64":
			return platformSpec{
				ManifestOS:         "linux",
				ManifestArch:       "arm64",
				NpmPlatformPackage: "@cursor/sdk-linux-arm64",
			}, nil
		case "amd64":
			return platformSpec{
				ManifestOS:         "linux",
				ManifestArch:       "amd64",
				NpmPlatformPackage: "@cursor/sdk-linux-x64",
			}, nil
		}
	case "windows":
		if runtime.GOARCH == "amd64" {
			return platformSpec{
				ManifestOS:         "win32",
				ManifestArch:       "amd64",
				NpmPlatformPackage: "@cursor/sdk-win32-x64",
			}, nil
		}
	}
	return platformSpec{}, fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
}
