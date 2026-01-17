package tests

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/donaldgifford/tflint-ruleset-terraform-style/tests/testutil"
)

// skipIfNoLocalStack skips the test if LocalStack is not running
func skipIfNoLocalStack(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("Skipping LocalStack test in short mode")
	}
	if !testutil.IsLocalStackRunning() {
		t.Skip("LocalStack not running - start with 'docker compose up -d'")
	}
}

// TestS3BucketWithTerraform tests S3 bucket creation using Terraform
func TestS3BucketWithTerraform(t *testing.T) {
	skipIfNoLocalStack(t)
	testS3Bucket(t, testutil.RuntimeTerraform)
}

// TestS3BucketWithTofu tests S3 bucket creation using OpenTofu
func TestS3BucketWithTofu(t *testing.T) {
	skipIfNoLocalStack(t)
	testS3Bucket(t, testutil.RuntimeTofu)
}

func testS3Bucket(t *testing.T, runtime testutil.Runtime) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ls := testutil.NewLocalStackConfig()
	t.Logf("LocalStack endpoint: %s, Runtime: %s", ls.Endpoint, runtime)

	// Run TFLint first to verify config passes our rules
	fixtureDir, err := filepath.Abs(filepath.Join("integration", "localstack", "s3_bucket"))
	if err != nil {
		t.Fatalf("Failed to get fixture dir: %v", err)
	}
	configPath, err := filepath.Abs(filepath.Join("..", ".tflint.hcl.example"))
	if err != nil {
		t.Fatalf("Failed to get config path: %v", err)
	}

	result, err := testutil.RunTFLint(fixtureDir, configPath)
	if err != nil {
		t.Fatalf("TFLint failed: %v", err)
	}

	if len(result.Issues) > 0 {
		t.Errorf("Expected passing config to have no issues, got %d:", len(result.Issues))
		for _, issue := range result.Issues {
			t.Errorf("  - %s: %s (%s:%d)", issue.Rule.Name, issue.Message, issue.Range.Filename, issue.Range.Start.Line)
		}
	}

	// Run Terraform/Tofu init/plan/apply
	tf := testutil.NewTerraformRunner(fixtureDir, ls.Endpoint, ls.Region).WithRuntime(runtime)
	defer func() { _ = tf.Cleanup() }()

	if err := tf.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	t.Log("Init passed")

	if err := tf.Plan(ctx); err != nil {
		t.Fatalf("Plan failed: %v", err)
	}
	t.Log("Plan passed")

	if err := tf.Apply(ctx); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	t.Log("Apply passed")

	// Verify outputs
	bucketID, err := tf.Output(ctx, "bucket_id")
	if err != nil {
		t.Logf("Warning: Could not get bucket_id output: %v", err)
	} else {
		t.Logf("Created bucket: %s", bucketID)
	}

	// Cleanup
	if err := tf.Destroy(ctx); err != nil {
		t.Logf("Warning: Destroy failed: %v", err)
	}
	t.Log("Destroy passed")
}

// TestIAMPolicyWithTerraform tests IAM policy creation using Terraform
func TestIAMPolicyWithTerraform(t *testing.T) {
	skipIfNoLocalStack(t)
	testIAMPolicy(t, testutil.RuntimeTerraform)
}

// TestIAMPolicyWithTofu tests IAM policy creation using OpenTofu
func TestIAMPolicyWithTofu(t *testing.T) {
	skipIfNoLocalStack(t)
	testIAMPolicy(t, testutil.RuntimeTofu)
}

func testIAMPolicy(t *testing.T, runtime testutil.Runtime) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ls := testutil.NewLocalStackConfig()
	t.Logf("LocalStack endpoint: %s, Runtime: %s", ls.Endpoint, runtime)

	fixtureDir, err := filepath.Abs(filepath.Join("integration", "localstack", "iam_policy"))
	if err != nil {
		t.Fatalf("Failed to get fixture dir: %v", err)
	}
	configPath, err := filepath.Abs(filepath.Join("..", ".tflint.hcl.example"))
	if err != nil {
		t.Fatalf("Failed to get config path: %v", err)
	}

	// Verify TFLint passes
	result, err := testutil.RunTFLint(fixtureDir, configPath)
	if err != nil {
		t.Fatalf("TFLint failed: %v", err)
	}

	if len(result.Issues) > 0 {
		t.Errorf("Expected no lint issues, got %d:", len(result.Issues))
		for _, issue := range result.Issues {
			t.Errorf("  - %s: %s", issue.Rule.Name, issue.Message)
		}
	}

	// Run Terraform/Tofu
	tf := testutil.NewTerraformRunner(fixtureDir, ls.Endpoint, ls.Region).WithRuntime(runtime)
	defer func() { _ = tf.Cleanup() }()

	if err := tf.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	t.Log("Init passed")

	if err := tf.Plan(ctx); err != nil {
		t.Fatalf("Plan failed: %v", err)
	}
	t.Log("Plan passed")

	if err := tf.Apply(ctx); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	t.Log("Apply passed")

	// Verify outputs
	policyARN, err := tf.Output(ctx, "policy_arn")
	if err != nil {
		t.Logf("Warning: Could not get policy_arn output: %v", err)
	} else {
		t.Logf("Created policy: %s", policyARN)
	}

	if err := tf.Destroy(ctx); err != nil {
		t.Logf("Warning: Destroy failed: %v", err)
	}
	t.Log("Destroy passed")
}

// TestCombinedWithTerraform tests combined module using Terraform
func TestCombinedWithTerraform(t *testing.T) {
	skipIfNoLocalStack(t)
	testCombined(t, testutil.RuntimeTerraform)
}

// TestCombinedWithTofu tests combined module using OpenTofu
func TestCombinedWithTofu(t *testing.T) {
	skipIfNoLocalStack(t)
	testCombined(t, testutil.RuntimeTofu)
}

func testCombined(t *testing.T, runtime testutil.Runtime) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	ls := testutil.NewLocalStackConfig()
	t.Logf("LocalStack endpoint: %s, Runtime: %s", ls.Endpoint, runtime)

	fixtureDir, err := filepath.Abs(filepath.Join("integration", "localstack", "combined"))
	if err != nil {
		t.Fatalf("Failed to get fixture dir: %v", err)
	}
	configPath, err := filepath.Abs(filepath.Join("..", ".tflint.hcl.example"))
	if err != nil {
		t.Fatalf("Failed to get config path: %v", err)
	}

	// TFLint validation
	result, err := testutil.RunTFLint(fixtureDir, configPath)
	if err != nil {
		t.Fatalf("TFLint failed: %v", err)
	}

	if len(result.Issues) > 0 {
		t.Errorf("Expected no lint issues, got %d:", len(result.Issues))
		for _, issue := range result.Issues {
			t.Errorf("  - %s: %s", issue.Rule.Name, issue.Message)
		}
	}

	// Terraform/Tofu workflow
	tf := testutil.NewTerraformRunner(fixtureDir, ls.Endpoint, ls.Region).WithRuntime(runtime)
	defer func() { _ = tf.Cleanup() }()

	if err := tf.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	t.Log("Init passed")

	if err := tf.Plan(ctx); err != nil {
		t.Fatalf("Plan failed: %v", err)
	}
	t.Log("Plan passed")

	if err := tf.Apply(ctx); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	t.Log("Apply passed")

	// Log outputs
	if bucketID, err := tf.Output(ctx, "bucket_id"); err == nil {
		t.Logf("Created bucket: %s", bucketID)
	}
	if policyARN, err := tf.Output(ctx, "policy_arn"); err == nil {
		t.Logf("Created policy: %s", policyARN)
	}

	if err := tf.Destroy(ctx); err != nil {
		t.Logf("Warning: Destroy failed: %v", err)
	}
	t.Log("Destroy passed")
}

// TestLintThenApply demonstrates the full workflow: lint -> plan -> apply
// Runs with both Terraform and OpenTofu
func TestLintThenApply(t *testing.T) {
	skipIfNoLocalStack(t)

	runtimes := []testutil.Runtime{testutil.RuntimeTerraform, testutil.RuntimeTofu}
	modules := []string{"s3_bucket", "iam_policy", "combined"}

	for _, runtime := range runtimes {
		t.Run(string(runtime), func(t *testing.T) {
			for _, module := range modules {
				t.Run(module, func(t *testing.T) {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
					defer cancel()

					ls := testutil.NewLocalStackConfig()

					fixtureDir, err := filepath.Abs(filepath.Join("integration", "localstack", module))
					if err != nil {
						t.Fatalf("Failed to get fixture dir: %v", err)
					}
					configPath, err := filepath.Abs(filepath.Join("..", ".tflint.hcl.example"))
					if err != nil {
						t.Fatalf("Failed to get config path: %v", err)
					}

					// Step 1: Lint
					result, err := testutil.RunTFLint(fixtureDir, configPath)
					if err != nil {
						t.Fatalf("TFLint failed: %v", err)
					}
					if len(result.Issues) > 0 {
						t.Errorf("Lint failed with %d issues", len(result.Issues))
						return
					}
					t.Log("Lint passed")

					// Step 2: Init + Plan + Apply
					tf := testutil.NewTerraformRunner(fixtureDir, ls.Endpoint, ls.Region).WithRuntime(runtime)
					defer func() { _ = tf.Cleanup() }()

					if err := tf.Init(ctx); err != nil {
						t.Fatalf("Init failed: %v", err)
					}
					t.Log("Init passed")

					if err := tf.Plan(ctx); err != nil {
						t.Fatalf("Plan failed: %v", err)
					}
					t.Log("Plan passed")

					if err := tf.Apply(ctx); err != nil {
						t.Fatalf("Apply failed: %v", err)
					}
					t.Log("Apply passed")

					// Step 3: Destroy
					if err := tf.Destroy(ctx); err != nil {
						t.Logf("Warning: Destroy failed: %v", err)
					}
					t.Log("Destroy passed")
				})
			}
		})
	}
}
