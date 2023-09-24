package acme

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/go-acme/lego/acme"
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

func NewACMERegistrationResource() resource.Resource {
	return &acmeRegistrationResource{}
}

type acmeRegistrationResource struct {
	serverUrl    string
	emailAddress string
}

type acmeRegistrationResourceModel struct {
	RegistrationUrl        types.String           `tfsdk:"registration_url"`
	AccountKeyPem          types.String           `tfsdk:"account_key_pem"`
	ExternalAccountBinding externalAccountBinding `tfsdk:"external_account_binding"`
}

type externalAccountBinding struct {
	KeyId      types.String `tfsdk:"key_id"`
	HmacBase64 types.String `tfsdk:"hmac_base64"`
}

func (r *acmeRegistrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registration"
}

func (r *acmeRegistrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"registration_url": schema.StringAttribute{
				Description: "ACME Registration URL.",
				Computed:    true,
			},
			"account_key_pem": schema.StringAttribute{
				Description: "ACME Account Key PEM",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Sensitive: true,
			},
			"external_account_binding": schema.SingleNestedAttribute{
				Description: "Domain to bind to instance domain.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"key_id": schema.StringAttribute{
						Description: "Key ID of EAB",
						Required:    true,
						Sensitive:   true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"hmac_base64": schema.StringAttribute{
						Description: "HMAC key of EAB in base 64 form",
						Required:    true,
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

func (r *acmeRegistrationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.serverUrl = req.ProviderData.(acmeProvider).serverUrl
	r.emailAddress = req.ProviderData.(acmeProvider).emailAddress
}

func (r *acmeRegistrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan *acmeRegistrationResourceModel
	getStateDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// register and agree to the TOS
	client, _, err := expandACMEClient(plan.AccountKeyPem.ValueString(), r.emailAddress, r.serverUrl, "", false)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to generate ACME client.",
			err.Error(),
		)
		return
	}

	var reg *registration.Resource
	// If EAB was enabled, register using EAB.
	registerAccount := func() error {
		reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: true,
			Kid:                  plan.ExternalAccountBinding.KeyId.ValueString(),
			HmacEncoded:          plan.ExternalAccountBinding.HmacBase64.ValueString(),
		})

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
	reconnectBackoff.MaxElapsedTime = 30 * time.Minute
	err = backoff.Retry(registerAccount, reconnectBackoff)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to register with EAB.",
			err.Error(),
		)
		return
	}

	plan.RegistrationUrl = types.StringValue(reg.URI)

	// Set state to fully populated data
	setStateDiags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *acmeRegistrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state *acmeRegistrationResourceModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// register and agree to the TOS
	_, user, err := expandACMEClient(state.AccountKeyPem.ValueString(), r.emailAddress, r.serverUrl, "", true)
	if err != nil {
		if regGone(err) {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"[API ERROR] Failed to generate ACME client.",
				err.Error(),
			)
			return
		}
	}

	state.RegistrationUrl = types.StringValue(user.Registration.URI)

	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *acmeRegistrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *acmeRegistrationResourceModel
	getPlanDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *acmeRegistrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *acmeRegistrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, _, err := expandACMEClient(state.AccountKeyPem.ValueString(), r.emailAddress, r.serverUrl, "", true)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to generate ACME client.",
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
	reconnectBackoff.MaxElapsedTime = 30 * time.Minute
	err = backoff.Retry(deleteRegistration, reconnectBackoff)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to delete ACME registration.",
			err.Error(),
		)
		return
	}
}

func regGone(err error) bool {
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
