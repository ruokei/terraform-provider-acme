module "gcp_storage" {
  for_each = { for idx, gcpstore in local.gcp_storage_values : idx => gcpstore }

  source = "./gcp-storage/"

  resource_name = local.mod_tmpl.resource_name
  ca            = each.value.ca
  region        = each.value.region
  app_tag       = "${var.module_info.brand}-${var.module_info.env}-acme-client"
}
