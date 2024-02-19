package main

import (
	"fmt"

	"go.nhat.io/python3"
)

func main() {
	sys := python3.MustImportModule("sys")
	version := sys.GetAttr("version_info")

	pyMajor := version.GetAttr("major")
	defer pyMajor.DecRef()

	pyMinor := version.GetAttr("minor")
	defer pyMinor.DecRef()

	pyReleaseLevel := version.GetAttr("releaselevel")
	defer pyReleaseLevel.DecRef()

	major := python3.AsInt(pyMajor)
	minor := python3.AsInt(pyMinor)
	releaseLevel := python3.AsString(pyReleaseLevel)

	fmt.Println("major:", major)
	fmt.Println("minor:", minor)
	fmt.Println("release level:", releaseLevel)

	// Output:
	// major: 3
	// minor: 11
	// release level: final
}
