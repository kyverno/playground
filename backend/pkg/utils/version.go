package utils

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

func ParseKubeVersion(kubeVersion string) (string, error) {
	if kubeVersion == "" {
		return "1.30", nil
	}
	version, err := semver.NewVersion(kubeVersion)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(version.Major(), ".", version.Minor()), nil
}
