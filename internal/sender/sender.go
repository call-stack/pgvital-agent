package sender

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pgvitals/pgvitals/agent/internal/collector"
)

type Sender struct {
	serverURL  string
	apiKey     string
	instanceID string
	httpClient *http.Client
}

func New(serverURL, apiKey, instanceID string) *Sender {
	return &Sender{
		serverURL:  serverURL,
		apiKey:     apiKey,
		instanceID: instanceID,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *Sender) SetInstanceID(id string) {
	s.instanceID = id
}

type RegisterRequest struct {
	APIKey    string `json:"api_key"`
	Hostname  string `json:"hostname"`
	PgVersion string `json:"pg_version"`
}

type RegisterResponse struct {
	InstanceID string `json:"instance_id"`
	OrgID      string `json:"org_id"`
}

func (s *Sender) Register(ctx context.Context, hostname, pgVersion string) (*RegisterResponse, error) {
	body, err := json.Marshal(RegisterRequest{
		APIKey:    s.apiKey,
		Hostname:  hostname,
		PgVersion: pgVersion,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal register: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		s.serverURL+"/api/v1/agent/register", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("register: server returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("register: decode response: %w", err)
	}
	return &result, nil
}

func (s *Sender) SendSnapshot(ctx context.Context, snap *collector.Snapshot) error {
	snap.InstanceID = s.instanceID

	body, err := json.Marshal(snap)
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(body); err != nil {
		return fmt.Errorf("gzip write: %w", err)
	}
	gz.Close()

	req, err := http.NewRequestWithContext(ctx, "POST",
		s.serverURL+"/api/v1/agent/snapshots", &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("X-Instance-ID", s.instanceID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
