locals {
  alioss_buckets = [
    for region in var.alicloud_oss_regions : "${var.resource_name}-${var.ca}-${region}"
  ]

  aws_s3s = [
    for region in var.aws_s3_regions : "${var.resource_name}-${var.ca}-${region}"
  ]

  gcp_storage = [
    for region in var.gcp_storage_regions : "${var.resource_name}-${var.ca}-${lower(region)}"
  ]
}
