package qdrant

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQcloudAuthTokenJSON(t *testing.T) {
	t.Parallel()

	tok, err := parseQcloudAuthTokenJSON([]byte(`{"access_token":"tok","token_type":"Bearer","scope":"manage"}`))
	require.NoError(t, err)
	assert.Equal(t, "tok", tok)

	_, err = parseQcloudAuthTokenJSON([]byte(`{}`))
	require.Error(t, err)

	_, err = parseQcloudAuthTokenJSON([]byte(``))
	require.Error(t, err)
}

func TestCLITokenSource_CachesUntilExpiry(t *testing.T) {
	t.Parallel()

	calls := 0
	exp := time.Now().Add(10 * time.Minute).Unix()
	token := fakeJWT(t, exp)

	src := newCLITokenSource("qcloud", "grpc.cloud.qdrant.io:443")
	src.run = func(_ context.Context, name string, args ...string) ([]byte, error) {
		calls++
		assert.Equal(t, "qcloud", name)
		assert.Equal(t, []string{"auth", "token", "--json", "--endpoint", "grpc.cloud.qdrant.io:443"}, args)
		return []byte(fmt.Sprintf(`{"access_token":%q,"token_type":"Bearer","scope":"manage"}`, token)), nil
	}

	got1, err := src.Token(context.Background())
	require.NoError(t, err)
	got2, err := src.Token(context.Background())
	require.NoError(t, err)

	assert.Equal(t, token, got1)
	assert.Equal(t, token, got2)
	assert.Equal(t, 1, calls)
}

func TestCLITokenSource_RefetchesNearExpiry(t *testing.T) {
	t.Parallel()

	calls := 0
	exp := time.Now().Add(30 * time.Second).Unix() // within cliTokenSkew
	token := fakeJWT(t, exp)

	src := newCLITokenSource("qcloud", "")
	src.run = func(_ context.Context, _ string, args ...string) ([]byte, error) {
		calls++
		assert.Equal(t, []string{"auth", "token", "--json"}, args)
		return []byte(fmt.Sprintf(`{"access_token":%q}`, token)), nil
	}

	_, err := src.Token(context.Background())
	require.NoError(t, err)
	_, err = src.Token(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, calls)
}

func TestProviderConfig_AuthorizationHeader(t *testing.T) {
	t.Parallel()

	apiCfg := &ProviderConfig{AuthMode: authModeAPIKey, ApiKey: "secret"}
	h, err := apiCfg.authorizationHeader(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "apikey secret", h)

	src := newCLITokenSource("qcloud", "")
	src.run = func(context.Context, string, ...string) ([]byte, error) {
		return []byte(`{"access_token":"oauth-tok"}`), nil
	}
	cliCfg := &ProviderConfig{AuthMode: authModeCLI, CLITokens: src}
	h, err = cliCfg.authorizationHeader(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Bearer oauth-tok", h)
}

func fakeJWT(t *testing.T, exp int64) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload, err := json.Marshal(map[string]int64{"exp": exp})
	require.NoError(t, err)
	body := base64.RawURLEncoding.EncodeToString(payload)
	return header + "." + body + ".sig"
}
