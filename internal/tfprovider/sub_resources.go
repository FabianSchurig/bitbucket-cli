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
