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
}

var subResources = []subResource{
	{
		TypeName:    "workspace-hooks",
		Description: "Manage webhooks for a Bitbucket workspace",
		sourceOps:   func() []OperationDef { return WorkspacesResourceGroup.AllOps },
	},
	{
		TypeName:    "default-reviewers",
		Description: "Manage default reviewers for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return PRResourceGroup.AllOps },
	},
	{
		TypeName:    "project-default-reviewers",
		Description: "Manage default reviewers for a Bitbucket project",
		sourceOps:   func() []OperationDef { return ProjectsResourceGroup.AllOps },
	},
	{
		TypeName:    "pipeline-variables",
		Description: "Manage pipeline variables for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "workspace-pipeline-variables",
		Description: "Manage pipeline variables for a Bitbucket workspace",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "deployment-variables",
		Description: "Manage deployment environment variables for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "repo-group-permissions",
		Description: "Manage explicit group permissions for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return ReposResourceGroup.AllOps },
	},
	{
		TypeName:    "repo-user-permissions",
		Description: "Manage explicit user permissions for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return ReposResourceGroup.AllOps },
	},
	{
		TypeName:    "project-group-permissions",
		Description: "Manage explicit group permissions for a Bitbucket project",
		sourceOps:   func() []OperationDef { return ProjectsResourceGroup.AllOps },
	},
	{
		TypeName:    "project-user-permissions",
		Description: "Manage explicit user permissions for a Bitbucket project",
		sourceOps:   func() []OperationDef { return ProjectsResourceGroup.AllOps },
	},
	{
		TypeName:    "repo-deploy-keys",
		Description: "Manage deploy keys for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return DeploymentsResourceGroup.AllOps },
	},
	{
		TypeName:    "project-deploy-keys",
		Description: "Manage deploy keys for a Bitbucket project",
		sourceOps:   func() []OperationDef { return DeploymentsResourceGroup.AllOps },
	},
	// ─── Wave 2: additional sub-resources ─────────────────────────────────────
	{
		TypeName:    "tags",
		Description: "Manage Git tags for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return RefsResourceGroup.AllOps },
	},
	{
		TypeName:    "pipeline-ssh-keys",
		Description: "Manage pipeline SSH key pair for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "pipeline-known-hosts",
		Description: "Manage pipeline SSH known hosts for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "pipeline-schedules",
		Description: "Manage pipeline schedules for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "pipeline-config",
		Description: "Manage pipeline configuration for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "ssh-keys",
		Description: "Manage SSH keys for a Bitbucket user",
		sourceOps:   func() []OperationDef { return UsersResourceGroup.AllOps },
	},
	{
		TypeName:    "current-user",
		Description: "Read the currently authenticated Bitbucket user",
		sourceOps:   func() []OperationDef { return UsersResourceGroup.AllOps },
	},
	{
		TypeName:    "forked-repository",
		Description: "Manage forked repositories in Bitbucket",
		sourceOps:   func() []OperationDef { return ReposResourceGroup.AllOps },
	},
	{
		TypeName:    "project-branching-model",
		Description: "Manage the branching model for a Bitbucket project",
		sourceOps:   func() []OperationDef { return BranchingModelResourceGroup.AllOps },
	},
	{
		TypeName:    "pipeline-oidc",
		Description: "Read pipeline OIDC configuration for a Bitbucket workspace",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "pipeline-oidc-keys",
		Description: "Read pipeline OIDC configuration keys for a Bitbucket workspace",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "workspace-members",
		Description: "List and read workspace members in Bitbucket",
		sourceOps:   func() []OperationDef { return WorkspacesResourceGroup.AllOps },
	},
	{
		TypeName:    "annotations",
		Description: "Manage report annotations for a Bitbucket commit",
		sourceOps:   func() []OperationDef { return ReportsResourceGroup.AllOps },
	},
	{
		TypeName:    "commit-file",
		Description: "Manage file contents via commits in a Bitbucket repository",
		sourceOps:   func() []OperationDef { return ReposResourceGroup.AllOps },
	},
	{
		TypeName:    "pr-comments",
		Description: "Manage comments on a Bitbucket pull request",
		sourceOps:   func() []OperationDef { return PRResourceGroup.AllOps },
	},
	{
		TypeName:    "issue-comments",
		Description: "Manage comments on a Bitbucket issue",
		sourceOps:   func() []OperationDef { return IssuesResourceGroup.AllOps },
	},
	// ─── Wave 3: DevSecOps-focused sub-resources ────────────────────────────
	{
		TypeName:    "repo-runners",
		Description: "Manage pipeline runners for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "workspace-runners",
		Description: "Manage pipeline runners for a Bitbucket workspace",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "pipeline-caches",
		Description: "Manage pipeline caches for a Bitbucket repository",
		sourceOps:   func() []OperationDef { return PipelinesResourceGroup.AllOps },
	},
	{
		TypeName:    "gpg-keys",
		Description: "Manage GPG keys for a Bitbucket user",
		sourceOps:   func() []OperationDef { return UsersResourceGroup.AllOps },
	},
	{
		TypeName:    "user-emails",
		Description: "Read email addresses for the current Bitbucket user",
		sourceOps:   func() []OperationDef { return UsersResourceGroup.AllOps },
	},
	{
		TypeName:    "hook-types",
		Description: "Read subscribable webhook event types in Bitbucket",
		sourceOps:   func() []OperationDef { return HooksResourceGroup.AllOps },
	},
	{
		TypeName:    "workspace-permissions",
		Description: "Read user permissions for a Bitbucket workspace",
		sourceOps:   func() []OperationDef { return WorkspacesResourceGroup.AllOps },
	},
	{
		TypeName:    "repo-settings",
		Description: "Manage repository settings inheritance in Bitbucket",
		sourceOps:   func() []OperationDef { return ReposResourceGroup.AllOps },
	},
}

func init() {
	for _, sr := range subResources {
		ops := sr.sourceOps()
		RegisterResourceGroup(ResourceGroup{
			TypeName:    sr.TypeName,
			Description: sr.Description,
			Ops:         MapCRUDOps(sr.TypeName, ops),
			AllOps:      ops,
		})
	}
}
