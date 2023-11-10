module "dns_challenge" {
  source = "./dns-challenge/"

  for_each = local.domains

  domain = each.key

  cloud_creds = {
    alicloud = {
      access_key = var.cloud_creds.alicloud.access_key
      secret_key = var.cloud_creds.alicloud.secret_key
    }
    aws = {
      region     = var.module_info.aws_region
      access_key = var.cloud_creds.aws.access_key
      secret_key = var.cloud_creds.aws.secret_key
    }
    cloudflare = {
      api_token = var.cloud_creds.cloudflare.api_token
    }
  }
}
