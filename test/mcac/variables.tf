##################################################
# Cloud Credentials
##################################################

variable "cloud_creds" {
  description = <<EOL
  Cloud credentials.
    - `alicloud`       : AliCloud credential.
      - `access_key` : AliCloud non-CN access key.
      - `secret_key` : AliCloud non-CN secret key.
    - `aws`            : AWS credential.
      - `access_key`   : AWS access key.
      - `secret_key`   : AWS secret key.
    - `cloudflare`     : Cloudflare credential.
      - `api_token`    : Cloudflare API Token.
    - `google`         : GCP credential.
      - `project`      : GCP access key.
      - `credentials`  : GCP secret key.
  EOL
  type = object({
    alicloud = object({
      access_key = string
      secret_key = string
    })
    aws = object({
      access_key = string
      secret_key = string
    })
    cloudflare = object({
      api_token = string
    })
    google = object({
      project     = string
      credentials = string
    })
  })
}

##################################################
# Module Template
##################################################

variable "module_info" {
  description = <<EOF
  Map of module's info details.
    - `brand`                  : Product brand.
    - `env`                    : Application environment.
    - `be_app_category`        : Backend application category.
    - `alicloud_region`        : AliCloud region.
    - `aws_region`             : AWS region.
    - `gcp_region`             : GCP region.
EOF
  type = object({
    brand                  = string
    env                    = string
    be_app_category        = string
    alicloud_region = string
    aws_region             = string
    gcp_region             = string
  })
}

variable "module_tmpl" {
  description = <<EOL
  Module template output format.
    - `resource_name`                                    : Resource name of module.
    - `multi_clouds_acme_client_vault_path`              : multi-clouds-acme-client vault path.
    - `vault_policy_multi_clouds_acme_client_vault_path` : Vault policy multi-clouds-acme-client vault path.
    - `vault_policy_ssl_cert_vault_path`                 : Vault policy ssl cert vault path.
EOL
  type = object({
    resource_name = optional(string, "{brand}-{env}-{be_app_category}-acme-client")
  })
}

##################################################
# General Variables
##################################################

variable "alicloud_oss_regions" {
  description = "Supported AliCloud OSS regions"
  type        = list(string)
}

variable "aws_s3_regions" {
  description = "Supported AWS S3 regions"
  type        = list(string)
}

variable "gcp_storage_regions" {
  description = "Supported GCP Storage regions"
  type        = list(string)
}

variable "acme_client" {
  description = <<EOL
  ACME Client Configuration.
    - `renew_before_days`      : Certificates are renewed before set amount of days.
    - `register_email_address` : Email address registered.
    - `ca_list`                : List of CA name used to create buckets.
    - `domain_sans`            : List of domain SANs
EOL
  type = object({
    renew_before_days      = number
    register_email_address = string
    ca_list                = list(string)
    domain_sans            = map(list(string))
  })
}

variable "extra_domains" {
  description = "Extra domains outside of 20twenty"
  type        = list(string)
  default     = [""]
}
