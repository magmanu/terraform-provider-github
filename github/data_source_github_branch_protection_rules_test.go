package github

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccGithubBranchProtectionRulesDataSource(t *testing.T) {

	t.Run("queries branch protection rules without error", func(t *testing.T) {
		randomID := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

		config := fmt.Sprintf(`
			resource "github_repository" "test" {
			  name = "tf-acc-test-%[1]s"
				auto_init = true
			}

			data "github_branch_protection_rules" "all" {
				repository = github_repository.test.name
			}
		`, randomID)

		check := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.#", "0"),
		)

		testCase := func(t *testing.T, mode string) {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { skipUnlessMode(t, mode) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: config,
						Check:  check,
					},
				},
			})
		}

		t.Run("with an anonymous account", func(t *testing.T) {
			t.Skip("anonymous account not supported for this operation")
		})

		t.Run("with an individual account", func(t *testing.T) {
			testCase(t, individual)
		})

		t.Run("with an organization account", func(t *testing.T) {
			testCase(t, organization)
		})

	})

	t.Run("queries branch protection", func(t *testing.T) {
		randomID := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

		config := fmt.Sprintf(`
			resource "github_repository" "test" {
			  name = "tf-acc-test-%[1]s"
				auto_init = true
			}

			resource "github_branch_protection" "protection" {
				repository_id = github_repository.test.id
			 	pattern = "main*"
				allows_deletions = false
				allows_force_pushes = false
				require_code_owner_reviews = true

			}
		`, randomID)

		config2 := config + `
			data "github_branch_protection_rules" "all" {
				repository = github_repository.test.name
			}
		`

		check := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.#", "1"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.pattern", "main*"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.allows_deletions", "true"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.allows_force_pushes", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.blocks_creations", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.dismisses_stale_reviews", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.is_admin_enforced", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.lock_allows_fetch_and_merge", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.lock_branch", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.require_last_push_approval", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.required_approving_review_count", "1"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.requires_approving_reviews", "true"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.require_code_owner_reviews", "true"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.requires_commit_signatures", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.requires_conversation_resolution", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.requires_linear_history", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.requires_deployments", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.requires_status_checks", "true"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.requires_strict_status_checks", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.restricts_pushes", "false"),
			resource.TestCheckResourceAttr("data.github_branch_protection_rules.all", "rules.0.restricts_review_dismissals", "false"),
		)

		testCase := func(t *testing.T, mode string) {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { skipUnlessMode(t, mode) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: config,
					},
					{
						Config: config2,
						Check:  check,
					},
				},
			})
		}

		t.Run("with an anonymous account", func(t *testing.T) {
			t.Skip("anonymous account not supported for this operation")
		})

		t.Run("with an individual account", func(t *testing.T) {
			testCase(t, individual)
		})

		t.Run("with an organization account", func(t *testing.T) {
			testCase(t, organization)
		})

	})
}
