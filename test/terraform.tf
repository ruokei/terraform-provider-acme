terraform {
  required_version = ">= 1.0, < 2.0"
  required_providers {
    zerossl = {
      source  = "toowoxx/zerossl"
      version = "~> 0.1"
    }
    acme = {
      source  = "example.local/myklst/acme"
      # version = "~> 0.0.1"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
}

provider "acme" {
  server_url = "https://acme.zerossl.com/v2/DV90"
}

resource "zerossl_eab_credentials" "zerossl_eab" {
  for_each = toset(local.domains)

  api_key = "cc4bfde128332a2797f9dd85f7ef6268"
}

resource "tls_private_key" "private_key" {
  for_each = toset(local.domains)

  algorithm = "RSA"
}

resource "acme_registration" "reg" {
  for_each = toset(local.domains)

  account_key_pem = tls_private_key.private_key[each.key].private_key_pem
  email_address   = "test@sige.la"

  external_account_binding {
    key_id      = zerossl_eab_credentials.zerossl_eab[each.key].kid
    hmac_base64 = zerossl_eab_credentials.zerossl_eab[each.key].hmac_key
  }
}

resource "acme_certificate" "certificate" {
  for_each = toset(local.domains)

  account_key_pem           = acme_registration.reg[each.key].account_key_pem
  min_days_remaining        = 120
  common_name               = each.key
  subject_alternative_names = []

  dns_challenge {
    provider = "route53"
    config = {
      "AWS_ACCESS_KEY_ID"     = "AKIA6KOWV4HJUCYIMKWI"
      "AWS_SECRET_ACCESS_KEY" = "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"
      "AWS_REGION"            = "ap-southeast-1"
    }
  }
}

locals {
  domains = ["sige-test3.com", "sige-test4.com", "sige-test5.com", "sige-test6.com", "sige-test7.com", "sige-test8.com"]
}
