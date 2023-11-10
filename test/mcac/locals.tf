locals {
  mod_tmpl = data.st-utilities_module_tmpl.template.module_tmpl

  alicloud_oss_values = flatten([
    for ca in var.acme_client.ca_list : [
      for region in var.alicloud_oss_regions : {
        ca     = ca
        region = region
      }
    ]
  ])

  aws_s3_values = flatten([
    for ca in var.acme_client.ca_list : [
      for region in var.aws_s3_regions : {
        ca     = ca
        region = region
      }
    ]
  ])

  gcp_storage_values = flatten([
    for ca in var.acme_client.ca_list : [
      for region in var.gcp_storage_regions : {
        ca     = ca
        region = region
      }
    ]
  ])

  customer_domains      = jsondecode(base64decode(data.external.docker_build.result.result))

  # Iterate over each domain that is in 20twenty
  customer_domains_list = distinct(flatten([
    for domain in local.customer_domains.domains : [domain.domain]
  ]))

  # Combine the domains with extra domains and then look for corresponding SANs based on input
  domains = {
    for domain in toset(compact(concat(var.extra_domains, local.customer_domains_list))) : domain => lookup(var.acme_client.domain_sans, domain, [])
  }
}
