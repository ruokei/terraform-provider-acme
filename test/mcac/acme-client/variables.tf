##################################################
# ACME Client Variables
##################################################

variable "register_email_address" {
  description = "ACME Registration Email Address"
  type        = string
}

variable "dns_challenge" {
  description = "DNS Challenge configuration."
  type = object({
    provider = string
    config   = map(string)
  })
}

variable "renew_before_days" {
  description = "Days clients renew before expiry"
  type        = number
}

variable "domain" {
  description = "Domain to Register"
  type        = string
}

variable "sans" {
  description = "List of Subject Alternate Names"
  type        = list(string)
}

##################################################
# Cloud Bucket Variables
##################################################

variable "resource_name" {
  description = "Resource name."
  type        = string
}

variable "alicloud_oss_regions" {
  description = "Alicloud OSS regions."
  type        = list(string)
}

variable "aws_s3_regions" {
  description = "AWS S3 regions."
  type        = list(string)
}

variable "gcp_storage_regions" {
  description = "GCP Storage regions."
  type        = list(string)
}

