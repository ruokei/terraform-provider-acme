variable "region" {
  description = <<EOF
Regions for creating AWS S3 buckets. Please refer
https://docs.aws.amazon.com/general/latest/gr/s3.html for details.
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

# ----------------------------------------------------------------------------------
# BECAUSE GOOGLE KSM KEYRING DESTROY BY SCHEDULER, SO WE DON'T USE ENCRYPTION FIRST.
# ----------------------------------------------------------------------------------
# variable "keyring_location" {
#   description = <<EOF
# Please refer https://cloud.google.com/kms/docs/locations for details.
# EOF
#   type        = string
# }
