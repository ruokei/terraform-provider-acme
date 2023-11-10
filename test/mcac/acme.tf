##################################################
# ACME Client
##################################################

module "acme_client" {

  source = "./acme-client/"

  providers = {
    acme.letsencrypt = acme.letsencrypt
    acme.zerossl     = acme.zerossl
    acme.gcp         = acme.gcp
  }

  for_each = local.domains

  register_email_address = var.acme_client.register_email_address

  dns_challenge = {
    provider = module.dns_challenge[each.key].dns_challenge.provider
    config   = module.dns_challenge[each.key].dns_challenge.config
  }

  renew_before_days = var.acme_client.renew_before_days
  domain            = each.key
  sans              = each.value

  resource_name        = local.mod_tmpl.resource_name
  alicloud_oss_regions = var.alicloud_oss_regions
  aws_s3_regions       = var.aws_s3_regions
  gcp_storage_regions  = var.gcp_storage_regions

  depends_on = [module.alicloud_oss, module.aws_s3, module.gcp_storage]
}
