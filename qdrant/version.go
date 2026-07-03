package qdrant

import (
	"fmt"
	"runtime/debug"
)

const (
	providerUserAgentPrefix = "terraform-provider-qdrant-cloud"
	providerVersionDev      = "dev"
)

// providerVersion is injected at release build time via ldflags.
var providerVersion = providerVersionDev

func providerUserAgent() string {
	version := providerVersion
	if version == "" || version == providerVersionDev {
		if v := versionFromBuildInfo(); v != "" {
			version = v
		}
	}
	if version == "" {
		version = providerVersionDev
	}
	return fmt.Sprintf("%s/%s", providerUserAgentPrefix, version)
}

func versionFromBuildInfo() string {
	const modulePath = "github.com/qdrant/terraform-provider-qdrant-cloud"

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	if bi.Main.Path == modulePath && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}

	return ""
}
