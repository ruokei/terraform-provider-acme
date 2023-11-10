output "alicloud_oss_buckets" {
  value = values({
    for idx, alioss in module.alicloud_oss :
    idx => alioss.bucket_name
  })
}

output "aws_s3_buckets" {
  value = values({
    for idx, awss3 in module.aws_s3 :
    idx => awss3.bucket_name
  })
}

output "gcp_storage_buckets" {
  value = values({
    for idx, gcpstorage in module.gcp_storage :
    idx => gcpstorage.bucket_name
  })
}
