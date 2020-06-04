package utils

import "testing"

func TestGetResourceFile(t *testing.T) {

}

func TestGetVersion(t *testing.T) {
	version := GetVersion("feature-refactor")
	t.Log(version)
}
