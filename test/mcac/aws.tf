module "aws_s3" {
  for_each = { for idx, awss3 in local.aws_s3_values : idx => awss3 }

  source = "./aws-s3/"

  resource_name = local.mod_tmpl.resource_name
  ca            = each.value.ca
  region        = each.value.region
  app_tag       = "${var.module_info.brand}-${var.module_info.env}-acme-client"
}
