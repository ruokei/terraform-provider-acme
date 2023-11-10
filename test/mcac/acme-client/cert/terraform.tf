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
      source  = "example.local/myklst/acme"
      # version = "~> 0.0.7"
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
