package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderUserAgent(t *testing.T) {
	old := providerVersion
	t.Cleanup(func() { providerVersion = old })

	providerVersion = "1.2.3"
	assert.Equal(t, "terraform-provider-qdrant-cloud/1.2.3", providerUserAgent())
}

func TestProviderUserAgent_DefaultDev(t *testing.T) {
	old := providerVersion
	t.Cleanup(func() { providerVersion = old })

	providerVersion = providerVersionDev
	assert.Equal(t, "terraform-provider-qdrant-cloud/dev", providerUserAgent())
}
