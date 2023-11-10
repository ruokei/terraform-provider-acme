output "alicloud_oss_buckets" {
  value = module.multi_clouds_acme_client.alicloud_oss_buckets
}

output "aws_s3_buckets" {
  value = module.multi_clouds_acme_client.aws_s3_buckets
}

output "gcp_storage_buckets" {
  value = module.multi_clouds_acme_client.gcp_storage_buckets
}
