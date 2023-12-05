package acme

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &acmeProvider{}
)

// New is a helper function to simplify provider server
func New() provider.Provider {
	return &acmeProvider{}
}

type acmeProvider struct{}

type acmeProviderModel struct {
	ServerUrl types.String `tfsdk:"server_url"`
}

// Metadata returns the provider type name.
func (p *acmeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "acme"
}

// Schema defines the provider-level schema for configuration data.
func (p *acmeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The ACME provider is used to update SSL certificates of CAs: ZeroSSL, LetsEncrypt, GCP. ",
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				Description: "ACME CA server URL.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a AliCloud API client for data sources and resources.
func (p *acmeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config acmeProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.ServerUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Unknown CA Server URL",
			"The provider cannot connect to the CA server as there is an unknown configuration value for the"+
				"CA server URL. Set the value statically in the configuration, or use the ACME_SERVER_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	var serverUrl string
	if !config.ServerUrl.IsNull() {
		serverUrl = config.ServerUrl.ValueString()
	} else {
		serverUrl = os.Getenv("ACME_SERVER_URL")
	}

	// If any of the expected configuration are missing, return
	// errors with provider-specific guidance.
	if serverUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Missing CA Server URL",
			"The provider cannot connect to the CA server as there is a "+
				"missing or empty value for the CA server URL. Set the "+
				"server URL value in the configuration or use the ACME_SERVER_URL "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = config.ServerUrl
}

func (p *acmeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *acmeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAcmeRegistrationResource,
		NewAcmeCertificateResource,
	}
}
