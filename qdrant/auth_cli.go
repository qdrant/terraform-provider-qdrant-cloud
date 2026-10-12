package qdrant

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	authModeAPIKey = "api_key"
	authModeCLI    = "cli"

	defaultQcloudBinary = "qcloud"
	cliTokenSkew       = 60 * time.Second
)

// cliTokenSource fetches short-lived OAuth access tokens via `qcloud auth token --json`.
// Tokens are cached until near JWT expiry (or for a short TTL if exp is unavailable).
type cliTokenSource struct {
	binary   string
	endpoint string

	mu     sync.Mutex
	token  string
	expiry time.Time

	// run is overridden in tests. Defaults to exec.CommandContext.
	run func(ctx context.Context, name string, args ...string) ([]byte, error)
}

func newCLITokenSource(binary, endpoint string) *cliTokenSource {
	if binary == "" {
		binary = defaultQcloudBinary
	}
	return &cliTokenSource{
		binary:   binary,
		endpoint: endpoint,
		run: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			cmd := exec.CommandContext(ctx, name, args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				msg := strings.TrimSpace(stderr.String())
				if msg == "" {
					msg = err.Error()
				}
				return nil, fmt.Errorf("qcloud auth token failed: %s", msg)
			}
			return stdout.Bytes(), nil
		},
	}
}

func (s *cliTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if s.token != "" && now.Before(s.expiry.Add(-cliTokenSkew)) {
		return s.token, nil
	}

	args := []string{"auth", "token", "--json"}
	if s.endpoint != "" {
		args = append(args, "--endpoint", s.endpoint)
	}

	out, err := s.run(ctx, s.binary, args...)
	if err != nil {
		return "", err
	}

	token, err := parseQcloudAuthTokenJSON(out)
	if err != nil {
		return "", err
	}

	s.token = token
	if exp, ok := jwtExpiry(token); ok {
		s.expiry = exp
	} else {
		// Auth0 access tokens are typically 1h; refresh conservatively without JWT exp.
		s.expiry = now.Add(30 * time.Minute)
	}
	return s.token, nil
}

type qcloudAuthTokenJSON struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func parseQcloudAuthTokenJSON(out []byte) (string, error) {
	out = bytes.TrimSpace(out)
	if len(out) == 0 {
		return "", fmt.Errorf("qcloud auth token returned empty output")
	}

	var payload qcloudAuthTokenJSON
	if err := json.Unmarshal(out, &payload); err != nil {
		return "", fmt.Errorf("parse qcloud auth token JSON: %w", err)
	}
	token := strings.TrimSpace(payload.AccessToken)
	if token == "" {
		return "", fmt.Errorf("qcloud auth token JSON missing access_token")
	}
	return token, nil
}

func jwtExpiry(token string) (time.Time, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return time.Time{}, false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, false
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil || claims.Exp == 0 {
		return time.Time{}, false
	}
	return time.Unix(claims.Exp, 0), true
}
