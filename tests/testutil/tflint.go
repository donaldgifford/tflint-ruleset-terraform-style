package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// TFLintResult represents the JSON output from tflint
type TFLintResult struct {
	Issues []TFLintIssue `json:"issues"`
	Errors []TFLintError `json:"errors"`
}

// TFLintIssue represents a single lint issue
type TFLintIssue struct {
	Rule    TFLintRule  `json:"rule"`
	Message string      `json:"message"`
	Range   TFLintRange `json:"range"`
}

// TFLintRule represents rule metadata
type TFLintRule struct {
	Name     string `json:"name"`
	Severity string `json:"severity"`
}

// TFLintRange represents the location of an issue
type TFLintRange struct {
	Filename string    `json:"filename"`
	Start    TFLintPos `json:"start"`
	End      TFLintPos `json:"end"`
}

// TFLintPos represents a position in a file
type TFLintPos struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// TFLintError represents a tflint error
type TFLintError struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// RunTFLint executes tflint on the given directory and returns results
func RunTFLint(dir string, configPath string) (*TFLintResult, error) {
	args := []string{
		"--format=json",
		"--force",
	}
	if configPath != "" {
		args = append(args, fmt.Sprintf("--config=%s", configPath))
	}

	cmd := exec.Command("tflint", args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// tflint returns non-zero exit code when issues are found, so we ignore the error
	_ = cmd.Run()

	// If no output, return empty result
	if stdout.Len() == 0 {
		return &TFLintResult{}, nil
	}

	var result TFLintResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse tflint output: %w\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	return &result, nil
}

// HasIssueForRule checks if results contain an issue for a specific rule
func (r *TFLintResult) HasIssueForRule(ruleName string) bool {
	for _, issue := range r.Issues {
		if issue.Rule.Name == ruleName {
			return true
		}
	}
	return false
}

// CountIssuesForRule counts issues for a specific rule
func (r *TFLintResult) CountIssuesForRule(ruleName string) int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Rule.Name == ruleName {
			count++
		}
	}
	return count
}

// GetIssuesForRule returns all issues for a specific rule
func (r *TFLintResult) GetIssuesForRule(ruleName string) []TFLintIssue {
	var issues []TFLintIssue
	for _, issue := range r.Issues {
		if issue.Rule.Name == ruleName {
			issues = append(issues, issue)
		}
	}
	return issues
}
