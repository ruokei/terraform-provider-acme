resource "alicloud_ram_user" "intl_ram_user" {
  name  = local.mod_tmpl.resource_name
  force = true
}

resource "alicloud_ram_access_key" "intl_access_key" {
  user_name = alicloud_ram_user.intl_ram_user.name
}

module "alicloud_oss" {
  for_each = { for idx, alioss in local.alicloud_oss_values : idx => alioss }

  source = "./ali-oss/"

  resource_name = local.mod_tmpl.resource_name
  ca            = each.value.ca
  region        = each.value.region
  app_tag       = "${var.module_info.brand}-${var.module_info.env}-acme-client"
}
