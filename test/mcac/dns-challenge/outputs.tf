output "dns_challenge" {
  value = {
    provider = local.modified_dns_provider
    config   = local.credential_keys[local.modified_dns_provider]
  }
}
