multi-clouds-acme-client
========================

Setup ACME client to store SSL certificates in multiple cloud buckets.

## Available Scenario for Testing

- basic

To perform basic module test. Working directory is `examples/basic`.

## Known Issues

- [ZeroSSL](https://github.com/cert-manager/cert-manager/issues/5867) API has rate limiting on and frequently hits error code 429 **(Too Many Requests)**.</br>
To perform terraform/terragrunt apply and destroy, use the `-parallelism=3` flag to limit the requests at once.

## Input Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| cloud\_creds | Cloud credentials.<br>    - `alicloud`       : AliCloud credential.<br>      - `non_cn`       : AliCloud non-CN credential.<br>        - `access_key` : AliCloud non-CN access key.<br>        - `secret_key` : AliCloud non-CN secret key.<br>      - `cn`           : AliCloud CN credential.<br>        - `access_key` : AliCloud CN access key.<br>        - `secret_key` : AliCloud CN secret key.<br>    - `aws`            : AWS credential.<br>      - `access_key`   : AWS access key.<br>      - `secret_key`   : AWS secret key.<br>    - `cloudflare`     : Cloudflare credential.<br>      - `api_token`    : Cloudflare API Token.<br>    - `google`         : GCP credential.<br>      - `project`      : GCP access key.<br>      - `credentials`  : GCP secret key. | <pre>object({<br>    alicloud = object({<br>      non_cn = object({<br>        access_key = string<br>        secret_key = string<br>      })<br>      cn = object({<br>        access_key = string<br>        secret_key = string<br>      })<br>    })<br>    aws = object({<br>      access_key = string<br>      secret_key = string<br>    })<br>    cloudflare = object({<br>      api_token = string<br>    })<br>    google = object({<br>      project     = string<br>      credentials = string<br>    })<br>  })</pre> | n/a | yes |
| module\_info | Map of module's info details.<br>    - `brand`                  : Product brand.<br>    - `env`                    : Application environment.<br>    - `be_app_category`        : Backend application category.<br>    - `alicloud_non_cn_region` : AliCloud Non-CN region.<br>    - `alicloud_cn_region`     : AliCloud CN region.<br>    - `aws_region`             : AWS region.<br>    - `gcp_region`             : GCP region. | <pre>object({<br>    brand                  = string<br>    env                    = string<br>    be_app_category        = string<br>    alicloud_non_cn_region = string<br>    alicloud_cn_region     = string<br>    aws_region             = string<br>    gcp_region             = string<br>  })</pre> | n/a | yes |
| module\_tmpl | Module template output format.<br>    - `resource_name`                                    : Resource name of module.<br>    - `resource_name_cn`                                 : Resource name of module for CN resources.<br>    - `multi_clouds_acme_client_vault_path`              : multi-clouds-acme-client vault path.<br>    - `vault_policy_multi_clouds_acme_client_vault_path` : Vault policy multi-clouds-acme-client vault path.<br>    - `vault_policy_ssl_cert_vault_path`                 : Vault policy ssl cert vault path. | <pre>object({<br>    resource_name                                    = optional(string, "{brand}-{env}-{be_app_category}-acme-client")<br>    resource_name_cn                                 = optional(string, "{brand}-{env}-{be_app_category}-cn-acme-client")<br>    multi_clouds_acme_client_vault_path              = optional(string, "devops/terraform/{brand}/{env}/{be_app_category}/multi-clouds-acme-client")<br>    vault_policy_multi_clouds_acme_client_vault_path = optional(string, "devops/data/terraform/{brand}/{env}/{be_app_category}/multi-clouds-acme-client")<br>    vault_policy_acme_client_vault_path              = optional(string, "devops/data/golang/acme-client")<br>    vault_policy_ssl_cert_vault_path                 = optional(string, "devops/data/ssl_cert/{brand}/{env}/+/*")<br>  })</pre> | n/a | yes |
| alicloud\_oss\_regions | Supported AliCloud OSS regions | `list(string)` | n/a | yes |
| aws\_s3\_regions | Supported AWS S3 regions | `list(string)` | n/a | yes |
| gcp\_storage\_regions | Supported GCP Storage regions | `list(string)` | n/a | yes |
| acme\_client | ACME Client Configuration.<br>    - `renew_before_days`      : Certificates are renewed before set amount of days.<br>    - `register_email_address` : Email address registered.<br>    - `ca_list`                : List of CA name used to create buckets.<br>    - `domains`              : List of domain information<br>      - `domain`             : Domain name<br>      - `sans`               : Subject alternative name | <pre>object({<br>    renew_before_days      = number<br>    register_email_address = string<br>    ca_list                = list(string)<br><br>    domains = list(object({<br>      domain = string<br>      sans   = list(string)<br>    }))<br>  })</pre> | n/a | yes |

## Output Variables

| Name | Description |
|------|-------------|
| alicloud\_cn\_user\_name | AliCloud CN account username, use for DNS challenge. |
| alicloud\_cn\_user\_access\_key | AliCloud CN account access key, use for DNS challenge. |
| alicloud\_cn\_user\_secret\_key | AliCloud CN account secret key, use for DNS challenge. |
| alicloud\_intl\_user\_name | AliCloud non-CN account username, use for OSS service. |
| alicloud\_intl\_user\_access\_key | AliCloud non-CN account access key, use for OSS service. |
| alicloud\_intl\_user\_secret\_key | AliCloud non-CN account secret key, use for OSS service. |
| aws\_user\_name | AWS account username, use for DNS challenge and S3 service. |
| aws\_user\_access\_key | AWS account access key, use for DNS challenge and S3 service. |
| aws\_user\_secret\_key | AWS account secret key, use for DNS challenge and S3 service. |
| cloudflare\_api\_token\_id | Cloudflare API token id, use for DNS challenge. |
| cloudflare\_api\_token | Cloudflare API token, use for DNS challenge. |
| gcloud\_storage\_user\_name | Google Cloud service account username, use for Cloud Storage service. |
| gcloud\_storage\_user\_public\_key | Google Cloud service account public key, use for Cloud Storage service. |

## How to add DNS providers

  1. Check the [ACME](https://registry.terraform.io/providers/vancluever/acme/latest/docs) provider, under **DNS Providers** to see if the DNS provider you are adding is present and supported.

  2. Dig to see the nameserver of the domain and validate if it is supported by the ACME provider. <br>Example output from: `dig sige-test.com`

  ```
  ;; AUTHORITY SECTION:
  sige-test.com.          1800    IN      SOA     clint.ns.cloudflare.com. dns.cloudflare.com. 2320128118 10000 2400 604800 1800
  ```

  3. Under `./dns-challenge/locals.tf`, add in the new provider within `supported_dns_providers` list.

  4. In cases where the nameserver of a domain does not contain the same wording as the **DNS Providers** in ACME provider<br>(`ns-1288.awsdns-33.org` corresponds to `route53` within the ACME provider), add the conversion into `provider_name_conversions`

  ```
    supported_dns_providers = ["alidns", "awsdns", "cloudflare"]
    provider_name_conversions = {
      "awsdns" = "route53"
      # Add more conversions as needed
    }
  ```

## How to add Certificate Authoritiies (CA)

  1. Go to `./terraform.tf`, add a corresponding provider block with an alias and corresponding URL for the CA.

  ```
    provider "acme" {
      alias      = "<ca-name>"
      server_url = "<ca-server-url>"
    }
  ```

  2. Go to `./acme.tf`, add a corresponding module for the specific CA and change the values in the `providers` block and `ca_provider` values.

  ```
    module "acme_client_<ca-name>" {
      providers = {
        acme = acme.<ca-name>
      }

      source = "./acme-client/"

      # Other configurations...

      resource_name        = local.mod_tmpl.resource_name
      ca_provider          = "<ca-name>"
      alicloud_oss_regions = var.alicloud_oss_regions
      aws_s3_regions       = var.aws_s3_regions
      gcp_storage_regions  = var.gcp_storage_regions
    }
  ```

## Additional Commands

  1. Terraform run apply with extra domains outside of 20twenty:<br>
  `terraform apply -var='extra_domains=["<extra-domains>"]'`

  2. Terraform run apply only for the targeted domains:<br>
  `terraform apply -target='module.multi_clouds_acme_client.module.acme_client["<domain>"]'`
