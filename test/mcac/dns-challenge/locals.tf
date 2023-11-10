locals {
  # DNS nameserver values
  supported_dns_providers = ["alidns", "awsdns", "cloudflare"]
  provider_name_conversions = {
    "awsdns" = "route53"
    # Add more conversions as needed
  }

  nameservers           = join(",", data.dns_ns_record_set.dns.nameservers)
  dns_provider          = compact([for provider in local.supported_dns_providers : can(regex(provider, local.nameservers)) ? provider : null])[0]
  modified_dns_provider = lookup(local.provider_name_conversions, local.dns_provider, local.dns_provider)

  credential_keys = {
    "cloudflare" = {
      "CLOUDFLARE_DNS_API_TOKEN"  = var.cloud_creds.cloudflare.api_token
      "CLOUDFLARE_ZONE_API_TOKEN" = var.cloud_creds.cloudflare.api_token
    }
    "route53" = {
      "AWS_REGION"            = var.cloud_creds.aws.region
      "AWS_ACCESS_KEY_ID"     = var.cloud_creds.aws.access_key
      "AWS_SECRET_ACCESS_KEY" = var.cloud_creds.aws.secret_key
    }
    "alidns" = {
      "ALICLOUD_ACCESS_KEY" = var.cloud_creds.alicloud.access_key
      "ALICLOUD_SECRET_KEY" = var.cloud_creds.alicloud.secret_key
    }
  }
}
