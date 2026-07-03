package qdrant

import (
	"fmt"
	"runtime/debug"
)

const providerUserAgentPrefix = "terraform-provider-qdrant-cloud"

// providerVersion is injected at release build time via ldflags.
var providerVersion = "dev"

func providerUserAgent() string {
	version := providerVersion
	if version == "" || version == "dev" {
		if v := versionFromBuildInfo(); v != "" {
			version = v
		}
	}
	if version == "" {
		version = "dev"
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
