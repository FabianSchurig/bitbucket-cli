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
	Username    types.String `tfsdk:"username"`
	AppPassword types.String `tfsdk:"app_password"`
	Token       types.String `tfsdk:"token"`
	BaseURL     types.String `tfsdk:"base_url"`
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
				Description: "Bitbucket username. Can also be set via BITBUCKET_USERNAME environment variable.",
				Optional:    true,
			},
			"app_password": schema.StringAttribute{
				Description: "Bitbucket app password. Can also be set via BITBUCKET_APP_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"token": schema.StringAttribute{
				Description: "Bitbucket OAuth2 access token. Can also be set via BITBUCKET_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Description: "Base URL for the Bitbucket API. Defaults to https://api.bitbucket.org/2.0. Can also be set via BITBUCKET_BASE_URL.",
				Optional:    true,
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

	// Set env vars from config if provided (config overrides env).
	if !config.Username.IsNull() {
		if err := os.Setenv("BITBUCKET_USERNAME", config.Username.ValueString()); err != nil {
			resp.Diagnostics.AddError("Failed to set BITBUCKET_USERNAME", err.Error())
			return
		}
	}
	if !config.AppPassword.IsNull() {
		if err := os.Setenv("BITBUCKET_APP_PASSWORD", config.AppPassword.ValueString()); err != nil {
			resp.Diagnostics.AddError("Failed to set BITBUCKET_APP_PASSWORD", err.Error())
			return
		}
	}
	if !config.Token.IsNull() {
		if err := os.Setenv("BITBUCKET_TOKEN", config.Token.ValueString()); err != nil {
			resp.Diagnostics.AddError("Failed to set BITBUCKET_TOKEN", err.Error())
			return
		}
	}
	if !config.BaseURL.IsNull() {
		if err := os.Setenv("BITBUCKET_BASE_URL", config.BaseURL.ValueString()); err != nil {
			resp.Diagnostics.AddError("Failed to set BITBUCKET_BASE_URL", err.Error())
			return
		}
	}

	// Validate that auth is configured.
	c, err := client.NewClient()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Bitbucket client",
			"Ensure BITBUCKET_USERNAME + BITBUCKET_APP_PASSWORD or BITBUCKET_TOKEN are set. Error: "+err.Error(),
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
)

// RegisterResourceGroup registers a resource group as both a Terraform resource
// and data source. Called by generated code at init time.
func RegisterResourceGroup(group ResourceGroup) {
	registeredResources = append(registeredResources, func() resource.Resource {
		return &GenericResource{group: group}
	})
	registeredDataSources = append(registeredDataSources, func() datasource.DataSource {
		return &GenericDataSource{group: group}
	})
}
