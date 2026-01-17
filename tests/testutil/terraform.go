package testutil

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Runtime specifies which IaC tool to use
type Runtime string

const (
	RuntimeTerraform Runtime = "terraform"
	RuntimeTofu      Runtime = "tofu"
)

// TerraformRunner manages Terraform/OpenTofu execution via Docker
type TerraformRunner struct {
	WorkDir     string // Local path to the fixture
	ContainerWD string // Path inside the container
	Runtime     Runtime
}

// NewTerraformRunner creates a new runner for the given directory
// workDir should be relative to the tests/ directory (e.g., "integration/localstack/s3_bucket")
func NewTerraformRunner(workDir, awsEndpoint, region string) *TerraformRunner {
	// Convert absolute path to relative path from tests/
	relPath := workDir
	if filepath.IsAbs(workDir) {
		// Extract the relative part after "tests/"
		if idx := strings.Index(workDir, "/tests/"); idx != -1 {
			relPath = workDir[idx+7:] // Skip "/tests/"
		}
	}

	return &TerraformRunner{
		WorkDir:     workDir,
		ContainerWD: "/workspace/" + relPath,
		Runtime:     RuntimeTerraform, // Default to Terraform
	}
}

// WithRuntime sets the runtime (terraform or tofu)
func (t *TerraformRunner) WithRuntime(r Runtime) *TerraformRunner {
	t.Runtime = r
	return t
}

// containerName returns the docker compose service name
func (t *TerraformRunner) containerName() string {
	return fmt.Sprintf("tflint-ruleset-terraform-style-%s-1", t.Runtime)
}

// runCommand executes a terraform/tofu command inside the container
func (t *TerraformRunner) runCommand(ctx context.Context, args ...string) error {
	// Build docker exec command
	dockerArgs := []string{
		"exec",
		"-w", t.ContainerWD,
		t.containerName(),
		string(t.Runtime),
	}
	dockerArgs = append(dockerArgs, args...)

	cmd := exec.CommandContext(ctx, "docker", dockerArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s failed: %s\n%s: %w",
			t.Runtime, args[0], stderr.String(), stdout.String(), err)
	}
	return nil
}

// runCommandWithOutput executes a command and returns stdout
func (t *TerraformRunner) runCommandWithOutput(ctx context.Context, args ...string) (string, error) {
	dockerArgs := []string{
		"exec",
		"-w", t.ContainerWD,
		t.containerName(),
		string(t.Runtime),
	}
	dockerArgs = append(dockerArgs, args...)

	cmd := exec.CommandContext(ctx, "docker", dockerArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s %s failed: %s: %w", t.Runtime, args[0], stderr.String(), err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Init runs terraform/tofu init
func (t *TerraformRunner) Init(ctx context.Context) error {
	return t.runCommand(ctx, "init", "-input=false")
}

// Plan runs terraform/tofu plan
func (t *TerraformRunner) Plan(ctx context.Context) error {
	return t.runCommand(ctx, "plan", "-input=false", "-out=tfplan")
}

// Apply runs terraform/tofu apply
func (t *TerraformRunner) Apply(ctx context.Context) error {
	return t.runCommand(ctx, "apply", "-auto-approve", "-input=false", "tfplan")
}

// Destroy runs terraform/tofu destroy
func (t *TerraformRunner) Destroy(ctx context.Context) error {
	return t.runCommand(ctx, "destroy", "-auto-approve", "-input=false")
}

// Output retrieves a terraform/tofu output value
func (t *TerraformRunner) Output(ctx context.Context, name string) (string, error) {
	return t.runCommandWithOutput(ctx, "output", "-raw", name)
}

// Cleanup removes terraform state files from the local directory
func (t *TerraformRunner) Cleanup() error {
	files := []string{
		"terraform.tfstate",
		"terraform.tfstate.backup",
		"tfplan",
		".terraform.lock.hcl",
	}
	for _, f := range files {
		_ = os.Remove(filepath.Join(t.WorkDir, f))
	}
	_ = os.RemoveAll(filepath.Join(t.WorkDir, ".terraform"))
	return nil
}
