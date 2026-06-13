package bridge

import "fmt"

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
