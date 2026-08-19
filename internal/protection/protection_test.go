package protection

import (
	"testing"

	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/config"
)

func TestBuildPlan(t *testing.T) {
	t.Parallel()
	cfg := config.Default()
	cfg.BranchProtection.Branch = "release/v1"
	plan := BuildPlan("owner/repo", cfg)
	if plan.Endpoint != "repos/owner/repo/branches/release%2Fv1/protection" {
		t.Fatalf("unexpected endpoint: %s", plan.Endpoint)
	}
	if plan.Payload["allow_force_pushes"] != false {
		t.Fatal("plan allows force pushes")
	}
}

func TestValidateRepository(t *testing.T) {
	t.Parallel()
	if err := ValidateRepository("owner/repo"); err != nil {
		t.Fatalf("valid repository rejected: %v", err)
	}
	for _, value := range []string{"repo", "owner/repo/extra", "../repo"} {
		if err := ValidateRepository(value); err == nil {
			t.Fatalf("invalid repository accepted: %q", value)
		}
	}
}
