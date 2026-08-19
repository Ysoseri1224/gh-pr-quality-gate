package protection

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/config"
)

var repositoryPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)

type Plan struct {
	Repository string                 `json:"repository"`
	Branch     string                 `json:"branch"`
	Endpoint   string                 `json:"endpoint"`
	Payload    map[string]interface{} `json:"payload"`
}

func BuildPlan(repository string, cfg config.Config) Plan {
	checks := make([]string, len(cfg.RequiredChecks))
	copy(checks, cfg.RequiredChecks)
	return Plan{
		Repository: repository,
		Branch:     cfg.BranchProtection.Branch,
		Endpoint:   fmt.Sprintf("repos/%s/branches/%s/protection", repository, url.PathEscape(cfg.BranchProtection.Branch)),
		Payload: map[string]interface{}{
			"required_status_checks": map[string]interface{}{"strict": true, "contexts": checks},
			"enforce_admins":         cfg.BranchProtection.EnforceAdmins,
			"required_pull_request_reviews": map[string]interface{}{
				"dismiss_stale_reviews":           cfg.BranchProtection.DismissStaleReviews,
				"required_approving_review_count": cfg.BranchProtection.RequiredApprovals,
				"require_last_push_approval":      cfg.BranchProtection.RequireLastPushApproval,
			},
			"restrictions":                     nil,
			"required_conversation_resolution": cfg.BranchProtection.RequireConversationResolution,
			"allow_force_pushes":               false,
			"allow_deletions":                  false,
		},
	}
}

func ValidateRepository(repository string) error {
	if !repositoryPattern.MatchString(repository) {
		return fmt.Errorf("repository must use owner/name form")
	}
	parts := strings.Split(repository, "/")
	if parts[0] == "." || parts[0] == ".." || parts[1] == "." || parts[1] == ".." {
		return fmt.Errorf("repository must use owner/name form")
	}
	return nil
}

func ResolveRepository(root string) (string, error) {
	cmd := exec.Command("gh", "repo", "view", "--json", "nameWithOwner", "--jq", ".nameWithOwner")
	cmd.Dir = root
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("resolve GitHub repository with gh: %w", err)
	}
	value := strings.TrimSpace(string(output))
	if value == "" {
		return "", fmt.Errorf("gh returned an empty repository name")
	}
	return value, nil
}

func Current(w io.Writer, root string, plan Plan) error {
	cmd := exec.Command("gh", "api", plan.Endpoint)
	cmd.Dir = root
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("read branch protection: %w", err)
	}
	return nil
}

func EnsureUnprotected(root string, plan Plan) error {
	cmd := exec.Command("gh", "api", "--silent", plan.Endpoint)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err == nil {
		return fmt.Errorf("branch protection already exists; refusing to replace settings that may be outside this tool's model")
	}
	if strings.Contains(string(output), "HTTP 404") {
		return nil
	}
	return fmt.Errorf("check existing branch protection: %s", strings.TrimSpace(string(output)))
}

func Apply(root string, plan Plan) error {
	file, err := os.CreateTemp("", "gh-pr-quality-gate-protection-*.json")
	if err != nil {
		return err
	}
	name := file.Name()
	defer os.Remove(name)
	if err := json.NewEncoder(file).Encode(plan.Payload); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	cmd := exec.Command("gh", "api", "--method", "PUT", plan.Endpoint, "--input", name)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("apply branch protection: %w", err)
	}
	return nil
}
