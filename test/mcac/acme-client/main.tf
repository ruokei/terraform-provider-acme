##################################################
# External Account Binding
##################################################

resource "st-gcp_acme_eab" "letsencrypt_eab" {}

resource "st-gcp_acme_eab" "gcp_eab" {}

resource "zerossl_eab_credentials" "zerossl_eab" {
  api_key = data.vault_generic_secret.zerossl.data["API_KEY"]
}

##################################################
# ACME Client
##################################################

module "acme_client_letsencrypt" {
  providers = {
    acme = acme.letsencrypt
  }

  source = "./cert/"

  register_email_address = var.register_email_address
  eab = {
    key_id      = st-gcp_acme_eab.letsencrypt_eab.key_id
    hmac_base64 = st-gcp_acme_eab.letsencrypt_eab.hmac_base64
  }

  dns_challenge = {
    provider = var.dns_challenge.provider
    config   = var.dns_challenge.config
  }

  renew_before_days = var.renew_before_days
  domain            = var.domain
  sans              = var.sans

  resource_name        = var.resource_name
  ca                   = "letsencrypt"
  alicloud_oss_regions = var.alicloud_oss_regions
  aws_s3_regions       = var.aws_s3_regions
  gcp_storage_regions  = var.gcp_storage_regions
}

module "acme_client_zerossl" {
  providers = {
    acme = acme.zerossl
  }

  source = "./cert/"

  register_email_address = var.register_email_address
  eab = {
    key_id      = zerossl_eab_credentials.zerossl_eab.kid
    hmac_base64 = zerossl_eab_credentials.zerossl_eab.hmac_key
  }

  dns_challenge = {
    provider = var.dns_challenge.provider
    config   = var.dns_challenge.config
  }

  renew_before_days = var.renew_before_days
  domain            = var.domain
  sans              = var.sans

  resource_name        = var.resource_name
  ca                   = "zerossl"
  alicloud_oss_regions = var.alicloud_oss_regions
  aws_s3_regions       = var.aws_s3_regions
  gcp_storage_regions  = var.gcp_storage_regions
}

module "acme_client_gcp" {
  providers = {
    acme = acme.gcp
  }

  source = "./cert/"

  register_email_address = var.register_email_address
  eab = {
    key_id      = st-gcp_acme_eab.gcp_eab.key_id
    hmac_base64 = st-gcp_acme_eab.gcp_eab.hmac_base64
  }

  dns_challenge = {
    provider = var.dns_challenge.provider
    config   = var.dns_challenge.config
  }

  renew_before_days = var.renew_before_days
  domain            = var.domain
  sans              = var.sans

  resource_name        = var.resource_name
  ca                   = "gcp"
  alicloud_oss_regions = var.alicloud_oss_regions
  aws_s3_regions       = var.aws_s3_regions
  gcp_storage_regions  = var.gcp_storage_regions
}
