terraform {
  required_providers {
    acme = {
      source = "example.local/myklst/acme"
    }
    zerossl = {
      source  = "toowoxx/zerossl"
      version = "~> 0.1"
    }
  }
}

provider "acme" {
  server_url = "https://acme.zerossl.com/v2/DV90"
}

resource "tls_private_key" "private_key" {
  algorithm = "RSA"
}

resource "zerossl_eab_credentials" "zerossl_eab" {
  api_key = "cc4bfde128332a2797f9dd85f7ef6268"
}

resource "acme_registration" "reg" {
  account_key_pem = tls_private_key.private_key.private_key_pem
  email_address   = "test@sige.la"

  external_account_binding {
    key_id      = zerossl_eab_credentials.zerossl_eab.kid
    hmac_base64 = zerossl_eab_credentials.zerossl_eab.hmac_key
  }
}

resource "acme_certificate" "certificate" {
  account_key_pem           = acme_registration.reg.account_key_pem
  email_address             = "test@sige.la"
  min_days_remaining        = 120
  common_name               = "sige-test3.com"
  subject_alternative_names = ["sige-test4.com"]

  dns_challenge {
    provider = "route53"
    config = {
      "AWS_REGION"            = "ap-southeast-1"
      "AWS_ACCESS_KEY_ID"     = "AKIA6KOWV4HJUCYIMKWI"
      "AWS_SECRET_ACCESS_KEY" = "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"
    }
  }
}
