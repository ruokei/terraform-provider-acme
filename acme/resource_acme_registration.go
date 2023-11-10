package acme

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-acme/lego/v4/acme"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &acmeRegistrationResource{}
	_ resource.ResourceWithConfigure = &acmeRegistrationResource{}
)

func NewAcmeRegistrationResource() resource.Resource {
	return &acmeRegistrationResource{}
}

type acmeRegistrationResource struct {
	ServerUrl types.String
}

type acmeRegistrationModel struct {
	RegistrationId         types.String           `tfsdk:"id"`
	AccountKeyPem          types.String           `tfsdk:"account_key_pem"`
	EmailAddress           types.String           `tfsdk:"email_address"`
	RegistrationUrl        types.String           `tfsdk:"registration_url"`
	ExternalAccountBinding externalAccountBinding `tfsdk:"external_account_binding"`
}

type externalAccountBinding struct {
	KeyId      types.String `tfsdk:"key_id"`
	HmacBase64 types.String `tfsdk:"hmac_base64"`
}

// Metadata returns the SSL binding resource name.
func (r *acmeRegistrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registration"
}

// Schema defines the schema for the SSL certificate binding resource.
func (r *acmeRegistrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Associates a domain with a SSL cert in Anti-DDoS website configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Registration ID.",
				Computed:    true,
			},
			"account_key_pem": schema.StringAttribute{
				Description: "Account Key PEM.",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email_address": schema.StringAttribute{
				Description: "Email address for CA.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"registration_url": schema.StringAttribute{
				Description: "CA Registration URL.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"external_account_binding": schema.SingleNestedBlock{
				Description: "External Account Binding configuration.",
				Attributes: map[string]schema.Attribute{
					"key_id": schema.StringAttribute{
						Description: "EAB Key ID.",
						Optional:    true,
						Sensitive:   true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"hmac_base64": schema.StringAttribute{
						Description: "HMAC Key in base64 format.",
						Optional:    true,
						Sensitive:   true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *acmeRegistrationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ServerUrl = types.StringValue(req.ProviderData.(string))
}

// Create a new SSL cert and domain binding
func (r *acmeRegistrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan *acmeRegistrationModel
	getStateDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, _, err := r.registrationClient(plan, r.ServerUrl.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] CREATE: Failed to create registration client",
			err.Error(),
		)
		return
	}

	var reg *registration.Resource

	// If EAB was enabled, register using EAB.
	registerAccount := func() error {
		if plan.ExternalAccountBinding.KeyId.ValueString() != "" && plan.ExternalAccountBinding.HmacBase64.ValueString() != "" {
			reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
				TermsOfServiceAgreed: true,
				Kid:                  plan.ExternalAccountBinding.KeyId.ValueString(),
				HmacEncoded:          plan.ExternalAccountBinding.HmacBase64.ValueString(),
			})
		} else {
			// Normal registration.
			reg, err = client.Registration.Register(registration.RegisterOptions{
				TermsOfServiceAgreed: true,
			})
		}
		if err != nil {
			if isAbleToRetry(err.Error()) {
				return err
			} else {
				return backoff.Permanent(err)
			}
		}
		return nil
	}
	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = DefaultMaxElapsedTime
	err = backoff.Retry(registerAccount, reconnectBackoff)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] CREATE: Failed to register with EAB.",
			err.Error(),
		)
		return
	}

	plan.RegistrationId = types.StringValue(reg.URI)

	// Set state to fully populated data
	setStateDiags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read web rules configuration for SSL cert and domain binding
func (r *acmeRegistrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state *acmeRegistrationModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, user, err := r.registrationClient(state, r.ServerUrl.ValueString(), true)
	if err != nil {
		if r.regGone(err) {
			resp.State.RemoveResource(ctx)
		}
		resp.Diagnostics.AddError(
			"[API ERROR] READ: Failed to create registration client",
			err.Error(),
		)
		return
	}

	state.RegistrationUrl = types.StringValue(user.Registration.URI)

	// Set state to fully populated data
	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update binds new SSL cert to domain and sets the updated Terraform state on success.
func (r *acmeRegistrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *acmeRegistrationModel
	var state *acmeRegistrationModel

	// Retrieve values from plan
	getPlanDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// SSL cert could not be unbinded, will always remain.
func (r *acmeRegistrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state *acmeRegistrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, _, err := r.registrationClient(state, r.ServerUrl.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] DELETE: Failed to create registration client",
			err.Error(),
		)
		return
	}

	deleteRegistration := func() error {
		err := client.Registration.DeleteRegistration()
		if err != nil {
			if isAbleToRetry(err.Error()) {
				return err
			} else {
				return backoff.Permanent(err)
			}
		}
		return nil
	}
	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = DefaultMaxElapsedTime
	err = backoff.Retry(deleteRegistration, reconnectBackoff)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] DELETE: Failed to delete registration",
			err.Error(),
		)
		return
	}

}

// registrationClient creates a connection to an ACME server from resource data,
// and also returns the user.
//
// If loadReg is supplied, the registration information is loaded in to the
// user's registration, if it exists - if the account cannot be resolved by the
// private key, then the appropriate error is returned.
func (r *acmeRegistrationResource) registrationClient(plan *acmeRegistrationModel, serverUrl string, loadReg bool) (*lego.Client, *acmeUser, error) {
	user := &acmeUser{
		key:   plan.AccountKeyPem.ValueString(),
		Email: plan.EmailAddress.ValueString(),
	}

	config := lego.NewConfig(user)
	config.CADirURL = serverUrl

	var client *lego.Client

	newClient := func() error {
		client, err := lego.NewClient(config)
		if err != nil {
			return err
		}

		// Populate user's registration resource if needed
		if loadReg {
			user.Registration, err = client.Registration.ResolveAccountByKey()
			if err != nil {
				if isAbleToRetry(err.Error()) {
					return err
				} else {
					return backoff.Permanent(err)
				}
			}
		}
		return nil
	}
	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = DefaultMaxElapsedTime
	err := backoff.Retry(newClient, reconnectBackoff)
	if err != nil {
		return nil, nil, err
	}

	return client, user, nil
}

func (r *acmeRegistrationResource) regGone(err error) bool {
	e, ok := err.(*acme.ProblemDetails)
	if !ok {
		return false
	}

	switch {
	case e.HTTPStatus == 400 && e.Type == "urn:ietf:params:acme:error:accountDoesNotExist":
		// As per RFC8555, see: no account exists when onlyReturnExisting
		// is set to true.
		return true

	case e.HTTPStatus == 403 && e.Type == "urn:ietf:params:acme:error:unauthorized":
		// Usually happens when the account has been deactivated. The URN
		// is a bit general for my liking, but it should be fine given
		// the specific nature of the request this error would be
		// returned for.
		return true
	}

	return false
}
