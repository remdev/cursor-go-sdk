package bridge

import "fmt"

// Test hooks exported for unit tests in the bridge_test package.

func EnsureInstalledForTestPlatform() (string, error) {
	plat, err := currentPlatform()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", plat.ManifestOS, plat.ManifestArch), nil
}

func NpmPlatformPackageForTest() (string, error) {
	plat, err := currentPlatform()
	if err != nil {
		return "", err
	}
	return plat.NpmPlatformPackage, nil
}

func ModuleBridgeDirForTest() string {
	return moduleBridge
}

func ValidBridgeRootForTest(root string) bool {
	return validBridgeRoot(root)
}
