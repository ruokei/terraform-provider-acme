package acme

import (
	"context"

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

type acmeProvider struct {
	serverUrl    string
	emailAddress string
}

type acmeProviderModel struct {
	ServerURL    types.String `tfsdk:"server_url"`
	EmailAddress types.String `tfsdk:"email_address"`
}

// Metadata returns the provider type name.
func (p *acmeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "acme"
}

// Schema defines the provider-level schema for configuration data.
func (p *acmeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The ACME client provider is used to automate the process of obtaining and managing digital certificates for securing websites and servers. " +
			"The provider needs to be configured with the proper configurations before it can be used.",
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				Description: "CA Server URL for ACME Client",
				Required:    true,
			},
			"email_address": schema.StringAttribute{
				Description: "Account Key PEM",
				Required:    true,
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
	if config.ServerURL.IsUnknown() || config.ServerURL.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Unknown/Missing CA Server URL",
			"The provider cannot ACME API client as there is an unknown or missing configuration value for the "+
				"ACME CA Server URL. Set the value statically in the configuration.",
		)
	}
	if config.EmailAddress.IsUnknown() || config.EmailAddress.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Unknown/Missing CA email address",
			"The provider cannot ACME API client as there is an unknown or missing configuration value for the "+
				"ACME CA email address. Set the value statically in the configuration.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	acmeProvder := acmeProvider{
		serverUrl:    config.ServerURL.ValueString(),
		emailAddress: config.EmailAddress.ValueString(),
	}
	resp.ResourceData = acmeProvder
}

func (p *acmeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *acmeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewACMERegistrationResource,
		NewACMECertificateResource,
	}
}
