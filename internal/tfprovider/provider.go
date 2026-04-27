// Package tfprovider implements the Terraform provider for Bitbucket Cloud.
//
// It exposes all Bitbucket API operations as Terraform resources and data sources,
// using a generic CRUD dispatch mechanism that reuses the same OperationDef metadata
// as the CLI and MCP server. Resource groups are auto-generated from the OpenAPI schema.
package tfprovider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
)

// Ensure the implementation satisfies the provider interface.
var _ provider.Provider = &BitbucketProvider{}

// BitbucketProvider implements the Terraform provider for Bitbucket Cloud.
type BitbucketProvider struct {
	version string
}

// BitbucketProviderModel describes the provider configuration.
type BitbucketProviderModel struct {
	Username          types.String `tfsdk:"username"`
	Token             types.String `tfsdk:"token"`
	BaseURL           types.String `tfsdk:"base_url"`
	CSRFToken         types.String `tfsdk:"csrf_token"`
	CloudSessionToken types.String `tfsdk:"cloud_session_token"`
}

// New creates a new BitbucketProvider instance.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BitbucketProvider{version: version}
	}
}

// Metadata returns the provider type name.
func (p *BitbucketProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bitbucket"
	resp.Version = p.version
}

// Schema defines the provider-level schema.
func (p *BitbucketProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for Bitbucket Cloud. Auto-generated from the Bitbucket OpenAPI spec.",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Description: "Bitbucket username (Atlassian account email). Can also be set via BITBUCKET_USERNAME environment variable. Required for Atlassian API tokens; not needed for workspace access tokens.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "Bitbucket API token (Atlassian API token or workspace access token). Can also be set via BITBUCKET_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Description: "Base URL for the Bitbucket API. Defaults to https://api.bitbucket.org/2.0. Can also be set via BITBUCKET_BASE_URL.",
				Optional:    true,
			},
			"csrf_token": schema.StringAttribute{
				Description: "CSRF token (value of the `csrftoken` browser cookie) used to authenticate against Bitbucket's internal API (`https://bitbucket.org/!api/internal/...`). Required for any resource that targets an internal endpoint, because those endpoints reject HTTP Basic Auth. Can also be set via the BITBUCKET_CSRF_TOKEN environment variable. Must be supplied together with `cloud_session_token`.",
				Optional:    true,
				Sensitive:   true,
			},
			"cloud_session_token": schema.StringAttribute{
				Description: "Cloud session token (value of the `cloud.session.token` browser cookie) used to authenticate against Bitbucket's internal API (`https://bitbucket.org/!api/internal/...`). Required for any resource that targets an internal endpoint. Can also be set via the BITBUCKET_CLOUD_SESSION_TOKEN environment variable. Must be supplied together with `csrf_token`.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure sets up the Bitbucket client from provider configuration.
func (p *BitbucketProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config BitbucketProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve values: config overrides env vars.
	username := os.Getenv("BITBUCKET_USERNAME")
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}
	token := os.Getenv("BITBUCKET_TOKEN")
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}
	baseURL := os.Getenv("BITBUCKET_BASE_URL")
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}
	csrfToken := os.Getenv("BITBUCKET_CSRF_TOKEN")
	if !config.CSRFToken.IsNull() {
		csrfToken = config.CSRFToken.ValueString()
	}
	cloudSessionToken := os.Getenv("BITBUCKET_CLOUD_SESSION_TOKEN")
	if !config.CloudSessionToken.IsNull() {
		cloudSessionToken = config.CloudSessionToken.ValueString()
	}

	c, err := client.NewClientWithConfig(username, token, baseURL, csrfToken, cloudSessionToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Bitbucket client",
			"Set BITBUCKET_TOKEN with a valid API token for the public REST API, "+
				"or set BITBUCKET_CSRF_TOKEN and BITBUCKET_CLOUD_SESSION_TOKEN to "+
				"target Bitbucket's internal API. Error: "+err.Error(),
		)
		return
	}

	// Share the client with resources and data sources.
	resp.DataSourceData = c
	resp.ResourceData = c
}

// Resources returns all registered Terraform resources.
func (p *BitbucketProvider) Resources(_ context.Context) []func() resource.Resource {
	return registeredResources
}

// DataSources returns all registered Terraform data sources.
func (p *BitbucketProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return registeredDataSources
}

// ─── Registration ─────────────────────────────────────────────────────────────

var (
	registeredResources   []func() resource.Resource
	registeredDataSources []func() datasource.DataSource
	registeredGroupList   []ResourceGroup
)

// RegisterResourceGroup registers a resource group as both a Terraform resource
// and data source. Called by generated code at init time.
func RegisterResourceGroup(group ResourceGroup) {
	registeredGroupList = append(registeredGroupList, group)
	registeredResources = append(registeredResources, func() resource.Resource {
		return &GenericResource{group: group}
	})
	registeredDataSources = append(registeredDataSources, func() datasource.DataSource {
		return &GenericDataSource{group: group}
	})
}

// RegisteredGroups returns all registered resource groups. Useful for code
// generators that need access to group metadata (params, operations, etc.).
func RegisteredGroups() []ResourceGroup {
	return registeredGroupList
}
