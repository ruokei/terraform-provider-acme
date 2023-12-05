package acme

import (
	"context"
	"crypto/x509"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/http/memcached"
	"github.com/go-acme/lego/v4/providers/http/s3"
	"github.com/go-acme/lego/v4/providers/http/webroot"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &acmeCertificateResource{}
	_ resource.ResourceWithConfigure        = &acmeCertificateResource{}
	_ resource.ResourceWithConfigValidators = &acmeCertificateResource{}
)

func NewAcmeCertificateResource() resource.Resource {
	return &acmeCertificateResource{}
}

type acmeCertificateResource struct {
	ServerUrl types.String
}

type acmeCertificateModel struct {
	CertificateId              types.String           `tfsdk:"id"`
	AccountKeyPem              types.String           `tfsdk:"account_key_pem"`
	EmailAddress               types.String           `tfsdk:"email_address"`
	CommonName                 types.String           `tfsdk:"common_name"`
	SubjectAlternativeNames    types.List             `tfsdk:"subject_alternative_names"`
	KeyType                    types.String           `tfsdk:"key_type"`
	CertificateRequestPem      types.String           `tfsdk:"certificate_request_pem"`
	MinDaysRemaining           types.Int64            `tfsdk:"min_days_remaining"`
	DnsChallenge               []*dnsChallenge         `tfsdk:"dns_challenge"`
	HttpChallenge              *httpChallenge          `tfsdk:"http_challenge"`
	HttpWebrootChallenge       *httpWebrootChallenge   `tfsdk:"http_webroot_challenge"`
	HttpMemcachedChallenge     *httpMemcachedChallenge `tfsdk:"http_memcached_challenge"`
	HttpS3Challenge            *httpS3Challenge        `tfsdk:"http_s3_challenge"`
	TlsChallenge               *tlsChallenge           `tfsdk:"tls_challenge"`
	PreCheckDelay              types.Int64            `tfsdk:"pre_check_delay"`
	RecursiveNameservers       types.List             `tfsdk:"recursive_nameservers"`
	DisableCompletePropagation types.Bool             `tfsdk:"disable_complete_propagation"`
	MustStaple                 types.Bool             `tfsdk:"must_staple"`
	PreferredChain             types.String           `tfsdk:"preferred_chain"`
	CertTimeout                types.Int64            `tfsdk:"cert_timeout"`
	CertificateUrl             types.String           `tfsdk:"certificate_url"`
	CertificateDomain          types.String           `tfsdk:"certificate_domain"`
	PrivateKeyPem              types.String           `tfsdk:"private_key_pem"`
	CertificatePem             types.String           `tfsdk:"certificate_pem"`
	IssuerPem                  types.String           `tfsdk:"issuer_pem"`
	CertificateP12             types.String           `tfsdk:"certificate_p12"`
	CertificateNotAfter        types.String           `tfsdk:"certificate_not_after"`
	CertificateP12Password     types.String           `tfsdk:"certificate_p12_password"`
	RevokeCertificateOnDestroy types.Bool             `tfsdk:"revoke_certificate_on_destroy"`
}

type dnsChallenge struct {
	Provider types.String `tfsdk:"provider"`
	Config   types.Map    `tfsdk:"config"`
}

type httpChallenge struct {
	Port        types.Int64  `tfsdk:"port"`
	ProxyHeader types.String `tfsdk:"proxy_header"`
}

type httpWebrootChallenge struct {
	Directory types.String `tfsdk:"directory"`
}

type httpMemcachedChallenge struct {
	Hosts types.List `tfsdk:"hosts"`
}

type httpS3Challenge struct {
	S3Bucket types.String `tfsdk:"s3_bucket"`
}

type tlsChallenge struct {
	Port types.Int64 `tfsdk:"port"`
}

// DNSProviderWrapper is a multi-provider wrapper to support multiple
// DNS challenges.
type DNSProviderWrapper struct {
	providers []challenge.ProviderTimeout
}

// CleanUp implements challenge.Provider.
func (*DNSProviderWrapper) CleanUp(domain string, token string, keyAuth string) error {
	panic("unimplemented")
}

// Present implements challenge.Provider.
func (*DNSProviderWrapper) Present(domain string, token string, keyAuth string) error {
	panic("unimplemented")
}

// NewDNSProviderWrapper returns an freshly initialized
// DNSProviderWrapper.
func NewDNSProviderWrapper() (*DNSProviderWrapper, error) {
	return &DNSProviderWrapper{}, nil
}

func (r *acmeCertificateResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("certificate_request_pem"),
			path.MatchRoot("common_name"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("certificate_request_pem"),
			path.MatchRoot("subject_alternative_names"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("certificate_request_pem"),
			path.MatchRoot("key_type"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("http_challenge"),
			path.MatchRoot("http_webroot_challenge"),
			path.MatchRoot("http_memcached_challenge"),
			path.MatchRoot("http_s3_challenge"),
		),
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("dns_challenge"),
			path.MatchRoot("http_challenge"),
			path.MatchRoot("http_webroot_challenge"),
			path.MatchRoot("http_memcached_challenge"),
			path.MatchRoot("http_s3_challenge"),
			path.MatchRoot("tls_challenge"),
		),
	}
}

// Metadata returns the SSL binding resource name.
func (r *acmeCertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

// Schema defines the schema for the SSL certificate binding resource.
func (r *acmeCertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Associates a domain with a SSL cert in Anti-DDoS website configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Certificate ID.",
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
				Description: "Registration Email Address.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"common_name": schema.StringAttribute{
				Description: "Domain name.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subject_alternative_names": schema.ListAttribute{
				Description: "Subject alternative names of domain.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"key_type": schema.StringAttribute{
				Description: "Key Type.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("P256", "P384", "2048", "4096", "8192"),
				},
			},
			"certificate_request_pem": schema.StringAttribute{
				Description: "Certificate Request PEM.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"min_days_remaining": schema.Int64Attribute{
				Description: "Minimum Days Remaining before certificate renews.",
				Optional:    true,
			},
			"pre_check_delay": schema.Int64Attribute{
				Description: "Pre-check delay time period.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"recursive_nameservers": schema.ListAttribute{
				Description: "Other nameservers.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"disable_complete_propagation": schema.BoolAttribute{
				Description: "Disable complete propagation.",
				Optional:    true,
			},
			"must_staple": schema.BoolAttribute{
				Description: "Must staple.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"preferred_chain": schema.StringAttribute{
				Description: "Preferred chain.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cert_timeout": schema.Int64Attribute{
				Description: "Certificate timeout.",
				Optional:    true,
			},
			"certificate_url": schema.StringAttribute{
				Description: "Certificate URL.",
				Computed:    true,
			},
			"certificate_domain": schema.StringAttribute{
				Description: "Certificate Domain.",
				Computed:    true,
			},
			"private_key_pem": schema.StringAttribute{
				Description: "Certificate Private Key PEM.",
				Computed:    true,
				Sensitive:   true,
			},
			"certificate_pem": schema.StringAttribute{
				Description: "Certificate PEM.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"issuer_pem": schema.StringAttribute{
				Description: "Certificate Issuer PEM.",
				Computed:    true,
			},
			"certificate_p12": schema.StringAttribute{
				Description: "Certificate P12.",
				Computed:    true,
				Sensitive:   true,
			},
			"certificate_not_after": schema.StringAttribute{
				Description: "Certificate not after.",
				Computed:    true,
			},
			"certificate_p12_password": schema.StringAttribute{
				Description: "Certificate P12 password.",
				Optional:    true,
				Sensitive:   true,
			},
			"revoke_certificate_on_destroy": schema.BoolAttribute{
				Description: "Whether to revoke certificate on resource destroy.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"dns_challenge": schema.ListNestedBlock{
				Description: "DNS Challenge.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"provider": schema.StringAttribute{
							Description: "DNS Provider.",
							Required:    true,
						},
						"config": schema.MapAttribute{
							Description: "DNS Provider Configuration.",
							Optional:    true,
							Sensitive:   true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"http_challenge": schema.SingleNestedBlock{
				Description: "HTTP Challenge.",
				Attributes: map[string]schema.Attribute{
					"port": schema.Int64Attribute{
						Description: "HTTP Port.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"proxy_header": schema.StringAttribute{
						Description: "HTTP Proxy Header.",
						Optional:    true,
					},
				},
			},
			"http_webroot_challenge": schema.SingleNestedBlock{
				Description: "HTTP Webroot Challenge.",
				Attributes: map[string]schema.Attribute{
					"directory": schema.StringAttribute{
						Description: "HTTP Webroot Directory.",
						Optional:    true,
					},
				},
			},
			"http_memcached_challenge": schema.SingleNestedBlock{
				Description: "HTTP Memcached Challenge.",
				Attributes: map[string]schema.Attribute{
					"hosts": schema.ListAttribute{
						Description: "HTTP Memcached hosts.",
						Optional:    true,
						ElementType: types.StringType,
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
						},
					},
				},
			},
			"http_s3_challenge": schema.SingleNestedBlock{
				Description: "HTTP s3 Challenge.",
				Attributes: map[string]schema.Attribute{
					"s3_bucket": schema.StringAttribute{
						Description: "s3 Challenge Bucket.",
						Optional:    true,
					},
				},
			},
			"tls_challenge": schema.SingleNestedBlock{
				Description: "TLS Challenge.",
				Attributes: map[string]schema.Attribute{
					"port": schema.Int64Attribute{
						Description: "TLS Port.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *acmeCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ServerUrl = req.ProviderData.(types.String)
}

// Create a new SSL cert and domain binding
func (r *acmeCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan *acmeCertificateModel
	getStateDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceUUID, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] CREATE: Unable to generate UUID.",
			err.Error(),
		)
		return
	}

	client, _, err := r.certificateClient(plan, r.ServerUrl.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] CREATE: Failed to create certificate client",
			err.Error(),
		)
		return
	}

	dnsCloser, err := r.setCertificateChallengeProviders(ctx, client, plan)
	defer dnsCloser()
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] CREATE: Failed to set certificate challenge providers",
			err.Error(),
		)
		return
	}

	var cert *certificate.Resource
	if !plan.CertificateRequestPem.IsNull() && !plan.CertificateRequestPem.IsUnknown() {
		var csr *x509.CertificateRequest
		csr, err = csrFromPEM([]byte(plan.CertificateRequestPem.ValueString()))
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] CREATE: Failed to convert CSR from PEM.",
				err.Error(),
			)
			return
		}

		var preferredChain string
		if !plan.PreferredChain.IsNull() && !plan.PreferredChain.IsUnknown() {
			preferredChain = plan.PreferredChain.ValueString()
		} else {
			preferredChain = ""
		}

		obtainCertCSR := func() error {
			cert, err = client.Certificate.ObtainForCSR(certificate.ObtainForCSRRequest{
				CSR:            csr,
				Bundle:         true,
				PreferredChain: preferredChain,
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
		reconnectBackoff.MaxElapsedTime = DefaultMaxElapsedTime
		err = backoff.Retry(obtainCertCSR, reconnectBackoff)
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] CREATE: Failed to Obtain cert for CSR.",
				err.Error(),
			)
			return
		}
	} else {
		domains := []string{plan.CommonName.ValueString()}
		if !plan.SubjectAlternativeNames.IsNull() && !plan.SubjectAlternativeNames.IsUnknown() {
			for _, x := range plan.SubjectAlternativeNames.Elements() {
				domains = append(domains, trimStringQuotes(x.String()))
			}
		}

		var mustStaple bool
		if !plan.MustStaple.IsNull() && !plan.MustStaple.IsUnknown() {
			mustStaple = plan.MustStaple.ValueBool()
		} else {
			mustStaple = false
		}

		var preferredChain string
		if !plan.PreferredChain.IsNull() && !plan.PreferredChain.IsUnknown() {
			preferredChain = plan.PreferredChain.ValueString()
		} else {
			preferredChain = ""
		}

		obtainCert := func() error {
			cert, err = client.Certificate.Obtain(certificate.ObtainRequest{
				Domains:        domains,
				Bundle:         true,
				MustStaple:     mustStaple,
				PreferredChain: preferredChain,
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
		reconnectBackoff.MaxElapsedTime = DefaultMaxElapsedTime
		err = backoff.Retry(obtainCert, reconnectBackoff)
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] CREATE: Failed to Obtain cert.",
				err.Error(),
			)
			return
		}
	}

	if len(cert.PrivateKey) > 0 {
		var certificateP12Password string
		if !plan.CertificateP12Password.IsNull() && !plan.CertificateP12Password.IsUnknown() {
			certificateP12Password = plan.CertificateP12Password.ValueString()
		} else {
			certificateP12Password = ""
		}

		pfxB64, err := bundleToPKCS12(cert.Certificate, cert.PrivateKey, certificateP12Password)
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] CREATE: Failed to bundle to PKCS12.",
				err.Error(),
			)
			return
		}
		plan.CertificateP12 = types.StringValue(string(pfxB64))
	} else {
		plan.CertificateP12 = types.StringValue(string(""))
	}

	issued, issuedNotAfter, issuer, err := splitPEMBundle(cert.Certificate)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] CREATE: Failed to split PEM bundle.",
			err.Error(),
		)
		return
	}

	plan.CertificateId = types.StringValue(resourceUUID)
	plan.CertificateUrl = types.StringValue(cert.CertURL)
	plan.CertificateDomain = types.StringValue(cert.Domain)
	plan.PrivateKeyPem = types.StringValue(string(cert.PrivateKey))
	plan.CertificatePem = types.StringValue(string(issued))
	plan.IssuerPem = types.StringValue(string(issuer))
	plan.CertificateNotAfter = types.StringValue(string(issuedNotAfter))

	// Set state to fully populated data
	setStateDiags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read web rules configuration for SSL cert and domain binding
func (r *acmeCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state *acmeCertificateModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	expired, err := resourceACMECertificateHasExpired(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] READ: Failed to check if certificate has expired.",
			err.Error(),
		)
		return
	}
	if expired {
		resp.State.RemoveResource(ctx)
	} else {
		// Try to recover the certificate from the ACME API.
		client, _, err := r.certificateClient(state, r.ServerUrl.ValueString(), true)
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] READ: Failed to create certificate client",
				err.Error(),
			)
			return
		}

		var cert *certificate.Resource
		getCert := func() error {
			cert, err = client.Certificate.Get(state.CertificateUrl.ValueString(), true)
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
		err = backoff.Retry(getCert, reconnectBackoff)
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] READ: Failed to get certificate",
				err.Error(),
			)
			return
		}

		if len(cert.PrivateKey) > 0 {
			var certificateP12Password string
			if !state.CertificateP12Password.IsNull() && !state.CertificateP12Password.IsUnknown() {
				certificateP12Password = state.CertificateP12Password.ValueString()
			} else {
				certificateP12Password = ""
			}

			pfxB64, err := bundleToPKCS12(cert.Certificate, cert.PrivateKey, certificateP12Password)
			if err != nil {
				resp.Diagnostics.AddError(
					"[API ERROR] READ: Failed to bundle to PKCS12.",
					err.Error(),
				)
				return
			}
			state.CertificateP12 = types.StringValue(string(pfxB64))
		} else {
			state.CertificateP12 = types.StringValue(string(""))
		}

		issued, issuedNotAfter, issuer, err := splitPEMBundle(cert.Certificate)
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] READ: Failed to split PEM bundle.",
				err.Error(),
			)
			return
		}

		state.CertificateUrl = types.StringValue(cert.CertURL)
		state.CertificateDomain = types.StringValue(cert.Domain)
		state.PrivateKeyPem = types.StringValue(string(cert.PrivateKey))
		state.CertificatePem = types.StringValue(string(issued))
		state.IssuerPem = types.StringValue(string(issuer))
		state.CertificateNotAfter = types.StringValue(string(issuedNotAfter))
	}

	// Set state to fully populated data
	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// resourceACMECertificateHasExpired checks the acme_certificate
// resource to see if it has expired.
func resourceACMECertificateHasExpired(state *acmeCertificateModel) (bool, error) {
	var mindays int64
	if !state.MinDaysRemaining.IsNull() && !state.MinDaysRemaining.IsUnknown() {
		mindays = state.MinDaysRemaining.ValueInt64()
	} else {
		mindays = 30
	}
	if mindays < 0 {
		log.Printf("[WARN] min_days_remaining is set to less than 0, certificate will never be renewed")
		return false, nil
	}

	cert := &certificate.Resource{
		Domain:  state.CertificateDomain.ValueString(),
		CertURL: state.CertificateUrl.ValueString(),
	}

	if !state.PrivateKeyPem.IsNull() && !state.PrivateKeyPem.IsUnknown() {
		cert.PrivateKey = []byte(state.PrivateKeyPem.ValueString())
	}

	if !state.CertificateRequestPem.IsNull() && !state.CertificateRequestPem.IsUnknown() {
		cert.CSR = []byte(state.CertificateRequestPem.ValueString())
	}

	remaining, err := certDaysRemaining(cert)
	if err != nil {
		return false, err
	}

	if int64(mindays) >= remaining {
		return true, nil
	}

	return false, nil
}

// Update binds new SSL cert to domain and sets the updated Terraform state on success.
func (r *acmeCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *acmeCertificateModel
	var state *acmeCertificateModel

	// Retrieve values from plan
	getPlanDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, _, err := r.certificateClient(state, r.ServerUrl.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] UPDATE: Failed to create certificate client",
			err.Error(),
		)
		return
	}

	cert := expandCertificateResource(state)

	dnsCloser, err := r.setCertificateChallengeProviders(ctx, client, state)
	defer dnsCloser()
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] UPDATE: Failed to set certificate challenge providers.",
			err.Error(),
		)
		return
	}

	var mustStaple bool
	if !plan.MustStaple.IsNull() && !plan.MustStaple.IsUnknown() {
		mustStaple = plan.MustStaple.ValueBool()
	} else {
		mustStaple = false
	}

	var preferredChain string
	if !plan.PreferredChain.IsNull() && !plan.PreferredChain.IsUnknown() {
		preferredChain = plan.PreferredChain.ValueString()
	} else {
		preferredChain = ""
	}

	var newCert *certificate.Resource
	renewCert := func() error {
		newCert, err = client.Certificate.Renew(*cert, true, mustStaple, preferredChain)
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
	err = backoff.Retry(renewCert, reconnectBackoff)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] UPDATE: Failed to renew certificate.",
			err.Error(),
		)
		return
	}

	state.CertificateP12Password = plan.CertificateP12Password
	if len(newCert.PrivateKey) > 0 {
		var certificateP12Password string
		if !state.CertificateP12Password.IsNull() && !state.CertificateP12Password.IsUnknown() {
			certificateP12Password = state.CertificateP12Password.ValueString()
		} else {
			certificateP12Password = ""
		}

		pfxB64, err := bundleToPKCS12(newCert.Certificate, newCert.PrivateKey, certificateP12Password)
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] UPDATE: Failed to bundle to PKCS12.",
				err.Error(),
			)
			return
		}
		state.CertificateP12 = types.StringValue(string(pfxB64))
	} else {
		state.CertificateP12 = types.StringValue(string(""))
	}

	issued, issuedNotAfter, issuer, err := splitPEMBundle(newCert.Certificate)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] UPDATE: Failed to split PEM bundle.",
			err.Error(),
		)
		return
	}

	state.CertificateUrl = types.StringValue(newCert.CertURL)
	state.CertificateDomain = types.StringValue(newCert.Domain)
	state.PrivateKeyPem = types.StringValue(string(newCert.PrivateKey))
	state.CertificatePem = types.StringValue(string(issued))
	state.IssuerPem = types.StringValue(string(issuer))
	state.CertificateNotAfter = types.StringValue(string(issuedNotAfter))

	// Set state to fully populated data
	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// SSL cert could not be unbinded, will always remain.
func (r *acmeCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state *acmeCertificateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var revokeCertificateOnDestroy bool
	if !state.RevokeCertificateOnDestroy.IsNull() && !state.RevokeCertificateOnDestroy.IsUnknown() {
		revokeCertificateOnDestroy = state.RevokeCertificateOnDestroy.ValueBool()
	} else {
		revokeCertificateOnDestroy = true
	}

	if !revokeCertificateOnDestroy {
		resp.Diagnostics.AddWarning(
			"[API ERROR] WARNING: Certificate not revoked.",
			"RevokeCertificateOnDestroy attribute set to false.",
		)
		return
	} else {
		client, _, err := r.certificateClient(state, r.ServerUrl.ValueString(), true)
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] DELETE: Failed to create certificate client.",
				err.Error(),
			)
			return
		}

		cert := expandCertificateResource(state)
		remaining, err := certSecondsRemaining(cert)
		if err != nil {
			resp.Diagnostics.AddError(
				"[API ERROR] DELETE: Failed to get certificate remaining seconds.",
				err.Error(),
			)
			return
		}

		if remaining >= 0 {
			revokeCert := func() error {
				err = client.Certificate.Revoke(cert.Certificate)
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
			err = backoff.Retry(revokeCert, reconnectBackoff)
			if err != nil {
				resp.Diagnostics.AddError(
					"[API ERROR] DELETE: Failed to revoke certificate.",
					err.Error(),
				)
				return
			}
		}
	}
}

func (r *acmeCertificateResource) certificateClient(plan *acmeCertificateModel, serverUrl string, loadReg bool) (*lego.Client, *acmeUser, error) {
	user := &acmeUser{
		key:   plan.AccountKeyPem.ValueString(),
		Email: plan.EmailAddress.ValueString(),
	}

	config := lego.NewConfig(user)
	config.CADirURL = serverUrl

	if !plan.KeyType.IsNull() && !plan.KeyType.IsUnknown() {
		config.Certificate.KeyType = certcrypto.KeyType(plan.KeyType.ValueString())
	} else {
		config.Certificate.KeyType = certcrypto.KeyType("2048")
	}

	if !plan.CertTimeout.IsNull() && !plan.CertTimeout.IsUnknown() {
		config.Certificate.Timeout = time.Duration(plan.CertTimeout.ValueInt64())
	} else {
		config.Certificate.Timeout = time.Duration(30)
	}

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

func (r *acmeCertificateResource) resourceACMECertificatePreCheckDelay(delay int) dns01.WrapPreCheckFunc {
	// Compute a reasonable interval for the delay, max delay 10
	// seconds, minimum 2.
	var interval int
	switch {
	case delay <= 10:
		interval = 2

	case delay <= 60:
		interval = 5

	default:
		interval = 10
	}

	return func(domain, fqdn, value string, orig dns01.PreCheckFunc) (bool, error) {
		stop, err := orig(fqdn, value)
		if stop && err == nil {
			// Run the delay. TODO: Eventually make this interruptible.
			var elapsed int
			end := time.After(time.Second * time.Duration(delay))
			for {
				select {
				case <-end:
					return true, nil
				default:
				}

				remaining := delay - elapsed
				if remaining < interval {
					// To honor the specified timeout, make our next interval the
					// time remaining. Minimum one second.
					interval = remaining
					if interval < 1 {
						interval = 1
					}
				}

				log.Printf("[DEBUG] [%s] acme: Waiting an additional %d second(s) for DNS record propagation.", domain, remaining)
				time.Sleep(time.Second * time.Duration(interval))
				elapsed += interval
			}
		}

		// A previous pre-check failed, return and exit.
		return stop, err
	}
}

func (r *acmeCertificateResource) setCertificateChallengeProviders(ctx context.Context, client *lego.Client, plan *acmeCertificateModel) (func(), error) {
	// DNSm
	dnsClosers := make([]func(), 0)
	dnsCloser := func() {
		for _, f := range dnsClosers {
			f()
		}
	}

	var nameservers []string
	if !plan.RecursiveNameservers.IsNull() && !plan.RecursiveNameservers.IsUnknown() {
		plan.RecursiveNameservers.ElementsAs(ctx, nameservers, false)
	}

	if plan.DnsChallenge != nil {
		dnsProvider, err := NewDNSProviderWrapper()
		if err != nil {
			return dnsCloser, err
		}

		for _, providerRaw := range plan.DnsChallenge {
			if p, closer, err := expandDNSChallenge(ctx, providerRaw, nameservers); err == nil {
				dnsProvider.providers = append(dnsProvider.providers, p)
				dnsClosers = append(dnsClosers, closer)
			} else {
				return dnsCloser, err
			}
		}

		setDns01Provider := func() error {
			dnsProvider, err := NewDNSProviderWrapper()
			if err != nil {
				return err
			}

			var preCheckDelay int64
			if plan.PreCheckDelay.IsNull() && plan.PreCheckDelay.IsUnknown() {
				preCheckDelay = plan.PreCheckDelay.ValueInt64()
			} else {
				preCheckDelay = 0
			}

			var disableCompletePropagation bool
			if plan.DisableCompletePropagation.IsNull() && plan.DisableCompletePropagation.IsUnknown() {
				disableCompletePropagation = plan.DisableCompletePropagation.ValueBool()
			} else {
				disableCompletePropagation = false
			}

			if err = client.Challenge.SetDNS01Provider(dnsProvider, r.expandDNSChallengeOptions(nameservers, disableCompletePropagation, preCheckDelay)...); err != nil {
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
		err = backoff.Retry(setDns01Provider, reconnectBackoff)
		if err != nil {
			return dnsCloser, err
		}
	}

	// HTTP (server)
	if reflect.DeepEqual(plan.HttpChallenge, httpChallenge{}) {
		var port int
		if !plan.HttpChallenge.Port.IsNull() && !plan.HttpChallenge.Port.IsUnknown() {
			port = int(plan.HttpChallenge.Port.ValueInt64())
		} else {
			port = 80
		}

		httpServerProvider := http01.NewProviderServer("", strconv.Itoa(port))
		if !plan.HttpChallenge.ProxyHeader.IsNull() && !plan.HttpChallenge.ProxyHeader.IsUnknown() {
			httpServerProvider.SetProxyHeader(plan.HttpChallenge.ProxyHeader.ValueString())
		}

		setHttp01Provider := func() error {
			if err := client.Challenge.SetHTTP01Provider(httpServerProvider); err != nil {
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
		err := backoff.Retry(setHttp01Provider, reconnectBackoff)
		if err != nil {
			return dnsCloser, err
		}
	}

	// HTTP (webroot)
	if reflect.DeepEqual(plan.HttpWebrootChallenge, httpWebrootChallenge{}) {
		httpWebrootProvider, err := webroot.NewHTTPProvider(plan.HttpWebrootChallenge.Directory.ValueString())
		if err != nil {
			return dnsCloser, err
		}

		setHttp01Provider := func() error {
			if err := client.Challenge.SetHTTP01Provider(httpWebrootProvider); err != nil {
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
		err = backoff.Retry(setHttp01Provider, reconnectBackoff)
		if err != nil {
			return dnsCloser, err
		}
	}

	// HTTP (memcached)
	if reflect.DeepEqual(plan.HttpMemcachedChallenge, httpMemcachedChallenge{}) {
		var hosts []string
		for _, host := range plan.HttpMemcachedChallenge.Hosts.Elements() {
			hosts = append(hosts, trimStringQuotes(host.String()))
		}

		httpMemcachedProvider, err := memcached.NewMemcachedProvider(hosts)
		if err != nil {
			return dnsCloser, err
		}

		setHttp01Provider := func() error {
			if err := client.Challenge.SetHTTP01Provider(httpMemcachedProvider); err != nil {
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
		err = backoff.Retry(setHttp01Provider, reconnectBackoff)
		if err != nil {
			return dnsCloser, err
		}
	}

	// HTTP (s3)
	if reflect.DeepEqual(plan.HttpS3Challenge, httpS3Challenge{}) {
		httpS3Provider, err := s3.NewHTTPProvider(plan.HttpS3Challenge.S3Bucket.ValueString())
		if err != nil {
			return dnsCloser, err
		}

		setHttp01Provider := func() error {
			if err := client.Challenge.SetHTTP01Provider(httpS3Provider); err != nil {
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
		err = backoff.Retry(setHttp01Provider, reconnectBackoff)
		if err != nil {
			return dnsCloser, err
		}
	}

	// TLS
	if reflect.DeepEqual(plan.TlsChallenge, tlsChallenge{}) {
		var port int
		if !plan.TlsChallenge.Port.IsNull() && !plan.TlsChallenge.Port.IsUnknown() {
			port = int(plan.TlsChallenge.Port.ValueInt64())
		} else {
			port = 443
		}

		tlsProvider := tlsalpn01.NewProviderServer("", strconv.Itoa(port))

		setTlsAlpn01Provider := func() error {
			if err := client.Challenge.SetTLSALPN01Provider(tlsProvider); err != nil {
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
		err := backoff.Retry(setTlsAlpn01Provider, reconnectBackoff)
		if err != nil {
			return dnsCloser, err
		}
	}

	return dnsCloser, nil
}

func (r *acmeCertificateResource) expandDNSChallengeOptions(nameservers []string, disableCompletePropagation bool, preCheckDelay int64) []dns01.ChallengeOption {
	var opts []dns01.ChallengeOption
	if len(nameservers) > 0 {
		opts = append(opts, dns01.AddRecursiveNameservers(nameservers))
	}

	if disableCompletePropagation {
		opts = append(opts, dns01.DisableCompletePropagationRequirement())
	}

	if preCheckDelay > 0 {
		opts = append(opts, dns01.WrapPreCheck(r.resourceACMECertificatePreCheckDelay(int(preCheckDelay))))
	}

	return opts
}

func expandCertificateResource(state *acmeCertificateModel) *certificate.Resource {
	cert := &certificate.Resource{
		Domain:  state.CertificateDomain.ValueString(),
		CertURL: state.CertificateUrl.ValueString(),
	}

	// Only populate the PrivateKey or CSR fields if we have them
	if !state.PrivateKeyPem.IsNull() && !state.PrivateKeyPem.IsUnknown() {
		cert.PrivateKey = []byte(state.PrivateKeyPem.ValueString())
	}

	if !state.CertificateRequestPem.IsNull() && !state.CertificateRequestPem.IsUnknown() {
		cert.PrivateKey = []byte(state.PrivateKeyPem.ValueString())
	}

	cert.Certificate = []byte(state.CertificatePem.ValueString() + state.IssuerPem.ValueString())

	return cert
}
