##################################################
# Create Certificate
##################################################

resource "tls_private_key" "private_key" {
  algorithm = "RSA"
}

# This null_resource is needed because submodule resources lifecycle block cannot refer to resources outside the submodule
# This null_resource will prompt a replace action when EAB changes
resource "null_resource" "eab_changed" {
  triggers = {
    "eab_key"        = var.eab.key_id
    "eab_hmacbase64" = var.eab.hmac_base64
  }
}

resource "acme_registration" "reg" {
  account_key_pem = tls_private_key.private_key.private_key_pem
  email_address   = var.register_email_address

  external_account_binding {
    key_id      = var.eab.key_id
    hmac_base64 = var.eab.hmac_base64
  }

  lifecycle {
    replace_triggered_by = [null_resource.eab_changed]
  }

  timeouts {
    create = "120m"
  }
}

resource "acme_certificate" "certificate" {
  account_key_pem           = acme_registration.reg.account_key_pem
  min_days_remaining        = var.renew_before_days
  common_name               = var.domain
  subject_alternative_names = var.sans

  dns_challenge {
    provider = var.dns_challenge.provider
    config   = var.dns_challenge.config
  }

  timeouts {
    create = "120m"
  }
}

##################################################
# OSS Bucket Objects
##################################################

resource "alicloud_oss_bucket_object" "privkey" {
  for_each = { for idx, bucket in local.alioss_buckets : idx => bucket }

  bucket  = each.value
  key     = "WILDCARD.${var.domain}/privkey.pem"
  content = acme_certificate.certificate.private_key_pem
}

resource "alicloud_oss_bucket_object" "cert" {
  for_each = { for idx, bucket in local.alioss_buckets : idx => bucket }

  bucket  = each.value
  key     = "WILDCARD.${var.domain}/cert.pem"
  content = acme_certificate.certificate.certificate_pem
}

resource "alicloud_oss_bucket_object" "chain" {
  for_each = { for idx, bucket in local.alioss_buckets : idx => bucket }

  bucket  = each.value
  key     = "WILDCARD.${var.domain}/chain.pem"
  content = acme_certificate.certificate.issuer_pem
}

resource "alicloud_oss_bucket_object" "fullchain" {
  for_each = { for idx, bucket in local.alioss_buckets : idx => bucket }

  bucket  = each.value
  key     = "WILDCARD.${var.domain}/fullchain.pem"
  content = <<-EOF
${acme_certificate.certificate.certificate_pem}
${acme_certificate.certificate.issuer_pem}
  EOF
}

resource "alicloud_oss_bucket_object" "full" {
  for_each = { for idx, bucket in local.alioss_buckets : idx => bucket }

  bucket  = each.value
  key     = "WILDCARD.${var.domain}/full.pem"
  content = <<-EOF
${acme_certificate.certificate.certificate_pem}
${acme_certificate.certificate.issuer_pem}
${acme_certificate.certificate.private_key_pem}
  EOF
}

##################################################
# AWS S3 Objects
##################################################

resource "aws_s3_bucket_object" "privkey" {
  for_each = { for idx, bucket in local.aws_s3s : idx => bucket }

  bucket  = each.value
  key     = "WILDCARD.${var.domain}/privkey.pem"
  content = acme_certificate.certificate.private_key_pem
}

resource "aws_s3_bucket_object" "cert" {
  for_each = { for idx, bucket in local.aws_s3s : idx => bucket }

  bucket  = each.value
  key     = "WILDCARD.${var.domain}/cert.pem"
  content = acme_certificate.certificate.certificate_pem
}

resource "aws_s3_bucket_object" "chain" {
  for_each = { for idx, bucket in local.aws_s3s : idx => bucket }

  bucket  = each.value
  key     = "WILDCARD.${var.domain}/chain.pem"
  content = acme_certificate.certificate.issuer_pem
}

resource "aws_s3_bucket_object" "fullchain" {
  for_each = { for idx, bucket in local.aws_s3s : idx => bucket }

  bucket  = each.value
  key     = "WILDCARD.${var.domain}/fullchain.pem"
  content = <<-EOF
${acme_certificate.certificate.certificate_pem}
${acme_certificate.certificate.issuer_pem}
  EOF
}

resource "aws_s3_bucket_object" "full" {
  for_each = { for idx, bucket in local.aws_s3s : idx => bucket }

  bucket  = each.value
  key     = "WILDCARD.${var.domain}/full.pem"
  content = <<-EOF
${acme_certificate.certificate.certificate_pem}
${acme_certificate.certificate.issuer_pem}
${acme_certificate.certificate.private_key_pem}
  EOF
}

##################################################
# GCP Storage Objects
##################################################

resource "google_storage_bucket_object" "privkey" {
  for_each = { for idx, bucket in local.gcp_storage : idx => bucket }

  bucket  = each.value
  name    = "WILDCARD.${var.domain}/privkey.pem"
  content = acme_certificate.certificate.private_key_pem
}

resource "google_storage_bucket_object" "cert" {
  for_each = { for idx, bucket in local.gcp_storage : idx => bucket }

  bucket  = each.value
  name    = "WILDCARD.${var.domain}/cert.pem"
  content = acme_certificate.certificate.certificate_pem
}

resource "google_storage_bucket_object" "chain" {
  for_each = { for idx, bucket in local.gcp_storage : idx => bucket }

  bucket  = each.value
  name    = "WILDCARD.${var.domain}/chain.pem"
  content = acme_certificate.certificate.issuer_pem
}

resource "google_storage_bucket_object" "fullchain" {
  for_each = { for idx, bucket in local.gcp_storage : idx => bucket }

  bucket  = each.value
  name    = "WILDCARD.${var.domain}/fullchain.pem"
  content = <<-EOF
${acme_certificate.certificate.certificate_pem}
${acme_certificate.certificate.issuer_pem}
  EOF
}

resource "google_storage_bucket_object" "full" {
  for_each = { for idx, bucket in local.gcp_storage : idx => bucket }

  bucket  = each.value
  name    = "WILDCARD.${var.domain}/full.pem"
  content = <<-EOF
${acme_certificate.certificate.certificate_pem}
${acme_certificate.certificate.issuer_pem}
${acme_certificate.certificate.private_key_pem}
  EOF
}
