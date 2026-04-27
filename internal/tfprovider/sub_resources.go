package tfprovider

// sub_resources.go registers additional Terraform resources derived from
// operations that already exist in generated resource groups. Each sub-resource
// selects a different set of CRUD operations via CRUDConfig, exposing
// Bitbucket sub-entities (e.g., workspace webhooks, default reviewers) as
// first-class Terraform resources.
//
// This file is hand-written. Add new entries when wiring additional
// sub-resources from the Bitbucket OpenAPI spec.

// subResource describes a sub-resource to register. The TypeName must have a
// corresponding entry in CRUDConfig.
type subResource struct {
	TypeName    string
	Description string
	// sourceOps returns the full list of operations from the parent group.
	// Using a func avoids init-order concerns with package-level vars.
	sourceOps func() []OperationDef
	// sourceCategory returns the parent group's Category string.
	sourceCategory func() string
}

func subResourceForGroup(typeName, description string, group *ResourceGroup) subResource {
	return subResource{
		TypeName:       typeName,
		Description:    description,
		sourceOps:      func() []OperationDef { return group.AllOps },
		sourceCategory: func() string { return group.Category },
	}
}

var subResources = []subResource{
	subResourceForGroup("workspace-hooks", "Manage webhooks for a Bitbucket workspace", &WorkspacesResourceGroup),
	subResourceForGroup("default-reviewers", "Manage default reviewers for a Bitbucket repository", &PRResourceGroup),
	subResourceForGroup("project-default-reviewers", "Manage default reviewers for a Bitbucket project", &ProjectsResourceGroup),
	subResourceForGroup("pipeline-variables", "Manage pipeline variables for a Bitbucket repository", &PipelinesResourceGroup),
	subResourceForGroup("workspace-pipeline-variables", "Manage pipeline variables for a Bitbucket workspace", &PipelinesResourceGroup),
	subResourceForGroup("deployment-variables", "Manage deployment environment variables for a Bitbucket repository", &PipelinesResourceGroup),
	subResourceForGroup("repo-group-permissions", "Manage explicit group permissions for a Bitbucket repository", &ReposResourceGroup),
	subResourceForGroup("repo-user-permissions", "Manage explicit user permissions for a Bitbucket repository", &ReposResourceGroup),
	subResourceForGroup("project-group-permissions", "Manage explicit group permissions for a Bitbucket project", &ProjectsResourceGroup),
	subResourceForGroup("project-user-permissions", "Manage explicit user permissions for a Bitbucket project", &ProjectsResourceGroup),
	subResourceForGroup("repo-deploy-keys", "Manage deploy keys for a Bitbucket repository", &DeploymentsResourceGroup),
	subResourceForGroup("project-deploy-keys", "Manage deploy keys for a Bitbucket project", &DeploymentsResourceGroup),
	// ─── Wave 2: additional sub-resources ─────────────────────────────────────
	subResourceForGroup("tags", "Manage Git tags for a Bitbucket repository", &RefsResourceGroup),
	subResourceForGroup("pipeline-ssh-keys", "Manage pipeline SSH key pair for a Bitbucket repository", &PipelinesResourceGroup),
	subResourceForGroup("pipeline-known-hosts", "Manage pipeline SSH known hosts for a Bitbucket repository", &PipelinesResourceGroup),
	subResourceForGroup("pipeline-schedules", "Manage pipeline schedules for a Bitbucket repository", &PipelinesResourceGroup),
	subResourceForGroup("pipeline-config", "Manage pipeline configuration for a Bitbucket repository", &PipelinesResourceGroup),
	subResourceForGroup("ssh-keys", "Manage SSH keys for a Bitbucket user", &UsersResourceGroup),
	subResourceForGroup("current-user", "Read the currently authenticated Bitbucket user", &UsersResourceGroup),
	subResourceForGroup("forked-repository", "Manage forked repositories in Bitbucket", &ReposResourceGroup),
	subResourceForGroup("project-branching-model", "Manage the branching model for a Bitbucket project", &BranchingModelResourceGroup),
	subResourceForGroup("pipeline-oidc", "Read pipeline OIDC configuration for a Bitbucket workspace", &PipelinesResourceGroup),
	subResourceForGroup("pipeline-oidc-keys", "Read pipeline OIDC configuration keys for a Bitbucket workspace", &PipelinesResourceGroup),
	subResourceForGroup("workspace-members", "List and read workspace members in Bitbucket", &WorkspacesResourceGroup),
	subResourceForGroup("annotations", "Manage report annotations for a Bitbucket commit", &ReportsResourceGroup),
	subResourceForGroup("commit-file", "Manage file contents via commits in a Bitbucket repository", &ReposResourceGroup),
	subResourceForGroup("pr-comments", "Manage comments on a Bitbucket pull request", &PRResourceGroup),
	subResourceForGroup("issue-comments", "Manage comments on a Bitbucket issue", &IssuesResourceGroup),
	// ─── Wave 3: DevSecOps-focused sub-resources ────────────────────────────
	subResourceForGroup("repo-runners", "Manage pipeline runners for a Bitbucket repository", &PipelinesResourceGroup),
	subResourceForGroup("workspace-runners", "Manage pipeline runners for a Bitbucket workspace", &PipelinesResourceGroup),
	subResourceForGroup("pipeline-caches", "Manage pipeline caches for a Bitbucket repository", &PipelinesResourceGroup),
	subResourceForGroup("gpg-keys", "Manage GPG keys for a Bitbucket user", &UsersResourceGroup),
	subResourceForGroup("user-emails", "Read email addresses for the current Bitbucket user", &UsersResourceGroup),
	subResourceForGroup("hook-types", "Read subscribable webhook event types in Bitbucket", &HooksResourceGroup),
	subResourceForGroup("workspace-permissions", "Read user permissions for a Bitbucket workspace", &WorkspacesResourceGroup),
	subResourceForGroup("repo-settings", "Manage repository settings inheritance in Bitbucket", &ReposResourceGroup),
	// Project-level branch restrictions (internal Bitbucket API) — exposed as
	// two sub-resources so users can manage glob-pattern and branching-model
	// rules independently. Both are powered by the same parent
	// ProjectBranchRestrictionsResourceGroup; CRUD wiring is in CRUDConfig.
	subResourceForGroup("project-branch-restrictions-by-pattern", "Manage project-level branch restrictions matched by glob pattern in Bitbucket (internal API)", &ProjectBranchRestrictionsResourceGroup),
	subResourceForGroup("project-branch-restrictions-by-branch-type", "Manage project-level branch restrictions matched by branching-model branch type in Bitbucket (internal API)", &ProjectBranchRestrictionsResourceGroup),
}

func init() {
	for _, sr := range subResources {
		ops := sr.sourceOps()
		RegisterResourceGroup(ResourceGroup{
			TypeName:    sr.TypeName,
			Description: sr.Description,
			Category:    sr.sourceCategory(),
			Ops:         MapCRUDOps(sr.TypeName, ops),
			AllOps:      ops,
		})
	}
}
