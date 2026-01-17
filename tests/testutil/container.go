package testutil

import (
	"fmt"
	"net/http"
	"time"
)

const (
	// LocalStackEndpoint is the default LocalStack endpoint
	LocalStackEndpoint = "http://localhost:4566"
	// LocalStackRegion is the default AWS region for LocalStack
	LocalStackRegion = "us-east-1"
)

// LocalStackConfig holds LocalStack connection info
type LocalStackConfig struct {
	Endpoint string
	Region   string
}

// NewLocalStackConfig returns the default LocalStack configuration
func NewLocalStackConfig() *LocalStackConfig {
	return &LocalStackConfig{
		Endpoint: LocalStackEndpoint,
		Region:   LocalStackRegion,
	}
}

// WaitForLocalStack waits for LocalStack to be ready
func WaitForLocalStack(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}

	for time.Now().Before(deadline) {
		resp, err := client.Get(LocalStackEndpoint + "/_localstack/health")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("LocalStack not ready after %v - ensure 'docker compose up -d' is running", timeout)
}

// IsLocalStackRunning checks if LocalStack is available
func IsLocalStackRunning() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(LocalStackEndpoint + "/_localstack/health")
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// AWSCredentials returns dummy credentials for LocalStack
func AWSCredentials() (accessKey, secretKey string) {
	return "test", "test"
}
