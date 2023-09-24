package acme

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &acmeCertificateResource{}
	_ resource.ResourceWithConfigure = &acmeCertificateResource{}
)

func NewACMECertificateResource() resource.Resource {
	return &acmeCertificateResource{}
}

type acmeCertificateResource struct {
	serverUrl    string
	emailAddress string
}

type acmeCertificateResourceModel struct {
	CertificateId           types.String `tfsdk:"certificate_id"`
	AccountKeyPem           types.String `tfsdk:"account_key_pem"`
	CommonName              types.String `tfsdk:"common_name"`
	SubjectAlternativeNames types.List   `tfsdk:"subject_alternative_names"`
	KeyType                 types.String `tfsdk:"key_type"`
	MinDaysRemaining        types.Int64  `tfsdk:"min_days_remaining"`
	CertificateUrl          types.String `tfsdk:"certificate_url"`
	CertificateDomain       types.String `tfsdk:"certificate_domain"`
	PrivateKeyPem           types.String `tfsdk:"private_key_pem"`
	CertificatePem          types.String `tfsdk:"certificate_pem"`
	IssuerPem               types.String `tfsdk:"issuer_pem"`
	CertificateP12          types.String `tfsdk:"certificate_p12"`
	CertificateNotAfter     types.String `tfsdk:"certificate_not_after"`
	DnsChallenge            types.List   `tfsdk:"dns_challenge"`
}

// type dnsChallenge struct {
// 	Provider types.String `tfsdk:"provider"`
// 	Config   types.Map    `tfsdk:"config"`
// }

func (r *acmeCertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (r *acmeCertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"certificate_id": schema.StringAttribute{
				Description: "ACME Certificate ID",
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
			"common_name": schema.StringAttribute{
				Description: "ACME certificate common name.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subject_alternative_names": schema.ListAttribute{
				Description: "ACME Certificate list of SANs.",
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"key_type": schema.StringAttribute{
				Description: "ACME Certificate Key Type. Valid values: P256, P384, 2048, 4096, 8192.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("2048"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("P256", "P384", "2048", "4096", "8192"),
				},
			},
			"min_days_remaining": schema.Int64Attribute{
				Description: "ACME minimum days remaining for certificate renewal.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(30),
			},
			"certificate_url": schema.StringAttribute{
				Description: "ACME Certificate URL.",
				Computed:    true,
			},
			"certificate_domain": schema.StringAttribute{
				Description: "ACME Certificate domain.",
				Computed:    true,
			},
			"private_key_pem": schema.StringAttribute{
				Description: "ACME Certificate private key PEM.",
				Computed:    true,
				Sensitive:   true,
			},
			"certificate_pem": schema.StringAttribute{
				Description: "ACME Certificate PEM.",
				Computed:    true,
			},
			"issuer_pem": schema.StringAttribute{
				Description: "ACME Issuer (chain) PEM.",
				Computed:    true,
			},
			"certificate_p12": schema.StringAttribute{
				Description: "ACME Certificate p12.",
				Computed:    true,
				Sensitive:   true,
			},
			"certificate_not_after": schema.StringAttribute{
				Description: "ACME Certificate expiration.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"dns_challenge": schema.ListNestedBlock{
				Description: "Certificate DNS Challenge",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"provider": schema.StringAttribute{
							Description: "DNS Provider",
							Required:    true,
						},
						"config": schema.MapAttribute{
							Description: "DNS Provider Configuration",
							Optional:    true,
							Sensitive:   true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *acmeCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.serverUrl = req.ProviderData.(acmeProvider).serverUrl
	r.emailAddress = req.ProviderData.(acmeProvider).emailAddress
}

func (r *acmeCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan *acmeCertificateResourceModel
	getStateDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to get plan",
			"can't get plan",
		)
		return
	}

	resourceUUID, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to generate UUID",
			err.Error(),
		)
		return
	}

	client, _, err := expandACMEClient(plan.AccountKeyPem.ValueString(), r.emailAddress, r.serverUrl, plan.KeyType.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to generate ACME client.",
			err.Error(),
		)
		return
	}

	var dns []dnsBlock
	for _, dnsChallenge := range plan.DnsChallenge.Elements() {
		var data dnsBlock
		json.Unmarshal([]byte(dnsChallenge.String()), &data)
		dns = append(dns, data)
	}

	dnsCloser, err := setCertificateChallengeProviders(client, dns)
	defer dnsCloser()
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Create: Failed to set certificate challenge providers",
			err.Error(),
		)
		return
	}

	var cert *certificate.Resource

	sans := []string{}
	domains := []string{plan.CommonName.ValueString()}

	_ = plan.SubjectAlternativeNames.ElementsAs(ctx, sans, false)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"[API ERROR] Create: Failed to convert SANs",
			"Terraform SANs input couldn't be converted.",
		)
		return
	}
	for _, san := range sans {
		if san != plan.CommonName.ValueString() {
			domains = append(domains, san)
		}
	}

	obtainCert := func() error {
		cert, err = client.Certificate.Obtain(certificate.ObtainRequest{
			Domains:        domains,
			Bundle:         true,
			MustStaple:     false,
			PreferredChain: "",
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
	err = backoff.Retry(obtainCert, reconnectBackoff)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Create: Failed to create certificate.",
			err.Error(),
		)
		return
	}

	plan.CertificateId = types.StringValue(resourceUUID)
	plan.CertificateUrl = types.StringValue(cert.CertURL)
	plan.CertificateDomain = types.StringValue(cert.Domain)
	plan.PrivateKeyPem = types.StringValue(string(cert.PrivateKey))

	issued, issuedNotAfter, issuer, err := splitPEMBundle(cert.Certificate)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Create: Failed to split PEM bundle.",
			err.Error(),
		)
		return
	}

	plan.CertificatePem = types.StringValue(string(issued))
	plan.IssuerPem = types.StringValue(string(issuer))
	plan.CertificateNotAfter = types.StringValue(issuedNotAfter)

	if len(cert.PrivateKey) > 0 {
		pfxB64, err := bundleToPKCS12(cert.Certificate, cert.PrivateKey, "")
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] Create: Failed to bundle certificate to PKCS12.",
				err.Error(),
			)
			return
		}

		plan.CertificateP12 = types.StringValue(string(pfxB64))
	} else {
		plan.CertificateP12 = types.StringValue("")
	}

	// Set state to fully populated data
	setStateDiags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"[API ERROR] Create: Failed to set plan into state",
			"can't set state.",
		)
		return
	}
}

func (r *acmeCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state *acmeCertificateResourceModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// client, _, err := expandACMEClient(state.AccountKeyPem.ValueString(), r.emailAddress, r.serverUrl, state.KeyType.ValueString(), true)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"[API ERROR] Read: Failed to generate ACME client.",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// var srcCR *certificate.Resource
	// getCert := func() error {
	// 	srcCR, err = client.Certificate.Get(state.CertificateUrl.ValueString(), true)
	// 	if err != nil {
	// 		// There are probably some cases that we will want to just drop
	// 		// the resource if there's been an issue, but seeing as this is
	// 		// mainly being used to recover for a bug that will be gone in
	// 		// 1.3.2, this will probably be rare. If we start relying on
	// 		// this behavior on a more general level, we may need to
	// 		// investigate this more. Just error on everything for now.
	// 		return err
	// 	}
	// 	return nil
	// }
	// reconnectBackoff := backoff.NewExponentialBackOff()
	// reconnectBackoff.MaxElapsedTime = 30 * time.Minute
	// err = backoff.Retry(getCert, reconnectBackoff)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"[API ERROR] Read: Failed to get certificate.",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// cert := &certificate.Resource{
	// 	Domain:     state.CertificateDomain.ValueString(),
	// 	CertURL:    state.CertificateUrl.ValueString(),
	// 	PrivateKey: []byte(state.PrivateKeyPem.ValueString()),
	// }
	// cert.Certificate = srcCR.Certificate

	// // Check if days remaining before cert expires is less than `min_days_remaining` set value
	// // and renew cert if its less than `min_days_remaining`
	// if state.MinDaysRemaining.ValueInt64() < 0 {
	// 	resp.Diagnostics.AddWarning(
	// 		"[API WARNING] Read: `min_days_remaining` is set to less than 0.",
	// 		"certificate will never be renewed",
	// 	)
	// } else {
	// 	remaining, err := certDaysRemaining(cert)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"[API ERROR] Read: Failed to calculate days remaining before cert expires.",
	// 			err.Error(),
	// 		)
	// 		return
	// 	}

	// 	if int64(state.MinDaysRemaining.ValueInt64()) >= remaining {
	// 		var providers []map[string]interface{}
	// 		respErr := state.DnsChallenge.ElementsAs(ctx, &providers, false)
	// 		if respErr.Errors() != nil {
	// 			resp.Diagnostics.Errors()
	// 			return
	// 		}

	// 		dnsCloser, err := setCertificateChallengeProviders(client, providers)
	// 		defer dnsCloser()
	// 		if err != nil {
	// 			resp.Diagnostics.AddError(
	// 				"[API ERROR] Read: Failed to set certificate challenge providers",
	// 				err.Error(),
	// 			)
	// 			return
	// 		}

	// 		var renewedCert *certificate.Resource

	// 		newCert := func() error {
	// 			renewedCert, err = client.Certificate.Renew(*cert, true, false, "")
	// 			if err != nil {
	// 				return err
	// 			}
	// 			return nil
	// 		}
	// 		reconnectBackoff := backoff.NewExponentialBackOff()
	// 		reconnectBackoff.MaxElapsedTime = 30 * time.Minute
	// 		err = backoff.Retry(newCert, reconnectBackoff)
	// 		if err != nil {
	// 			resp.Diagnostics.AddError(
	// 				"[API ERROR] Read: Failed to renew certificate",
	// 				err.Error(),
	// 			)
	// 			return
	// 		}

	// 		state.CertificateUrl = types.StringValue(renewedCert.CertURL)
	// 		state.CertificateDomain = types.StringValue(renewedCert.Domain)
	// 		state.PrivateKeyPem = types.StringValue(string(renewedCert.PrivateKey))

	// 		issued, issuedNotAfter, issuer, err := splitPEMBundle(renewedCert.Certificate)
	// 		if err != nil {
	// 			resp.Diagnostics.AddError(
	// 				"[API ERROR] Read: Failed to split PEM bundle.",
	// 				err.Error(),
	// 			)
	// 			return
	// 		}

	// 		state.CertificatePem = types.StringValue(string(issued))
	// 		state.IssuerPem = types.StringValue(string(issuer))
	// 		state.CertificateNotAfter = types.StringValue(issuedNotAfter)

	// 		if len(cert.PrivateKey) > 0 {
	// 			pfxB64, err := bundleToPKCS12(renewedCert.Certificate, renewedCert.PrivateKey, "")
	// 			if err != nil {
	// 				resp.Diagnostics.AddError(
	// 					"[API ERROR] Read: Failed to bundle certificate to PKCS12.",
	// 					err.Error(),
	// 				)
	// 				return
	// 			}
	// 			state.CertificateP12 = types.StringValue(string(pfxB64))
	// 		} else {
	// 			state.CertificateP12 = types.StringValue("")
	// 		}

	// 	}

	// }

	// state.CertificateUrl = types.StringValue(cert.CertURL)
	// state.CertificateDomain = types.StringValue(cert.Domain)
	// state.PrivateKeyPem = types.StringValue(string(cert.PrivateKey))

	// issued, issuedNotAfter, issuer, err := splitPEMBundle(cert.Certificate)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"[API ERROR] Read: Failed to split PEM bundle.",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// state.CertificatePem = types.StringValue(string(issued))
	// state.IssuerPem = types.StringValue(string(issuer))
	// state.CertificateNotAfter = types.StringValue(issuedNotAfter)

	// if len(cert.PrivateKey) > 0 {
	// 	pfxB64, err := bundleToPKCS12(cert.Certificate, cert.PrivateKey, "")
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"[API ERROR] Read: Failed to bundle certificate to PKCS12.",
	// 			err.Error(),
	// 		)
	// 		return
	// 	}
	// 	state.CertificateP12 = types.StringValue(string(pfxB64))
	// } else {
	// 	state.CertificateP12 = types.StringValue("")
	// }

	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *acmeCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *acmeCertificateResourceModel
	getPlanDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// client, _, err := expandACMEClient(state.AccountKeyPem.ValueString(), r.emailAddress, r.serverUrl, state.KeyType.ValueString(), true)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"[API ERROR] Update: Failed to generate ACME client.",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// cert := &certificate.Resource{
	// 	Domain:     state.CertificateDomain.ValueString(),
	// 	CertURL:    state.CertificateUrl.ValueString(),
	// 	PrivateKey: []byte(state.PrivateKeyPem.ValueString()),
	// }
	// cert.Certificate = []byte(state.CertificatePem.ValueString() + state.IssuerPem.ValueString())

	// var providers []map[string]interface{}
	// respErr := state.DnsChallenge.ElementsAs(ctx, &providers, false)
	// if respErr.Errors() != nil {
	// 	resp.Diagnostics.Errors()
	// 	return
	// }

	// dnsCloser, err := setCertificateChallengeProviders(client, providers)
	// defer dnsCloser()
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"[API ERROR] Update: Failed to set certificate challenge providers",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// var renewedCert *certificate.Resource

	// newCert := func() error {
	// 	renewedCert, err = client.Certificate.Renew(*cert, true, false, "")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }
	// reconnectBackoff := backoff.NewExponentialBackOff()
	// reconnectBackoff.MaxElapsedTime = 30 * time.Minute
	// err = backoff.Retry(newCert, reconnectBackoff)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"[API ERROR] Update: Failed to renew certificate",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// state.CertificateUrl = types.StringValue(renewedCert.CertURL)
	// state.CertificateDomain = types.StringValue(renewedCert.Domain)
	// state.PrivateKeyPem = types.StringValue(string(renewedCert.PrivateKey))

	// issued, issuedNotAfter, issuer, err := splitPEMBundle(renewedCert.Certificate)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"[API ERROR] Update: Failed to split PEM bundle.",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// state.CertificatePem = types.StringValue(string(issued))
	// state.IssuerPem = types.StringValue(string(issuer))
	// state.CertificateNotAfter = types.StringValue(issuedNotAfter)

	// if len(cert.PrivateKey) > 0 {
	// 	pfxB64, err := bundleToPKCS12(renewedCert.Certificate, renewedCert.PrivateKey, "")
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"[API ERROR] Update: Failed to bundle certificate to PKCS12.",
	// 			err.Error(),
	// 		)
	// 		return
	// 	}
	// 	state.CertificateP12 = types.StringValue(string(pfxB64))
	// } else {
	// 	state.CertificateP12 = types.StringValue("")
	// }

	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *acmeCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *acmeCertificateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// client, _, err := expandACMEClient(state.AccountKeyPem.ValueString(), r.emailAddress, r.serverUrl, state.KeyType.ValueString(), true)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"[API ERROR] Delete: Failed to generate ACME client.",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// cert := &certificate.Resource{
	// 	Domain:     state.CertificateDomain.ValueString(),
	// 	CertURL:    state.CertificateUrl.ValueString(),
	// 	PrivateKey: []byte(state.PrivateKeyPem.ValueString()),
	// }
	// cert.Certificate = []byte(state.CertificatePem.ValueString() + state.IssuerPem.ValueString())
	// remaining, err := certSecondsRemaining(cert)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"[API ERROR] Delete: Unable to get days remaining for certificate.",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// if remaining >= 0 {
	// 	revokeCert := func() error {
	// 		err = client.Certificate.Revoke(cert.Certificate)
	// 		if err != nil {
	// 			if isAbleToRetry(err.Error()) {
	// 				return err
	// 			} else {
	// 				return backoff.Permanent(err)
	// 			}
	// 		}
	// 		return nil
	// 	}
	// 	reconnectBackoff := backoff.NewExponentialBackOff()
	// 	reconnectBackoff.MaxElapsedTime = 30 * time.Minute
	// 	err = backoff.Retry(revokeCert, reconnectBackoff)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"[API ERROR] Delete: Unable to revoke certificate.",
	// 			err.Error(),
	// 		)
	// 		return
	// 	}
	// }
}
