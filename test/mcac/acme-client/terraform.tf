terraform {
  required_version = ">= 1.0, < 2.0"
  required_providers {
    st-gcp = {
      source  = "myklst/st-gcp"
      version = "~> 0.1"
    }
    zerossl = {
      source  = "toowoxx/zerossl"
      version = "0.1.1"
    }
    acme = {
      source                = "example.local/myklst/acme"
      configuration_aliases = [acme.letsencrypt, acme.zerossl, acme.gcp]
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
    alicloud = {
      source  = "aliyun/alicloud"
      version = "~> 1.0"
    }
  }
}
