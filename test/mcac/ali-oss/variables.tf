variable "region" {
  description = <<EOF
Regions for creating AliCloud OSS buckets. Please refer
https://www.alibabacloud.com/help/doc-detail/31837.htm for details.
EOF
  type        = string
}

variable "resource_name" {
  description = "Resource name for module."
  type        = string
}

variable "app_tag" {
  description = "App tag value."
  type        = string
}

variable "ca" {
  description = "Certificate authority name."
  type        = string
}
