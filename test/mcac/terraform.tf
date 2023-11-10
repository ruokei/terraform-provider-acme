terraform {
  required_version = ">= 1.0, < 2.0"
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "~> 1.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
    st-utilities = {
      source  = "myklst/st-utilities"
      version = "~> 0.1"
    }
    vault = {
      source  = "hashicorp/vault"
      version = "~> 3.0"
    }
    st-gcp = {
      source  = "myklst/st-gcp"
      version = "~> 0.1"
    }
    zerossl = {
      source  = "toowoxx/zerossl"
      version = "~> 0.1"
    }
    acme = {
      source  = "example.local/myklst/acme"
      # version = "~> 0.0.7"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
}

provider "alicloud" {
  region     = var.module_info.alicloud_region
  access_key = var.cloud_creds.alicloud.access_key
  secret_key = var.cloud_creds.alicloud.secret_key
}

provider "aws" {
  region     = var.module_info.aws_region
  access_key = var.cloud_creds.aws.access_key
  secret_key = var.cloud_creds.aws.secret_key
}

provider "google" {
  region      = var.module_info.gcp_region
  project     = var.cloud_creds.google.project
  credentials = var.cloud_creds.google.credentials
}

provider "st-gcp" {
  project     = var.cloud_creds.google.project
  credentials = var.cloud_creds.google.credentials
}

provider "cloudflare" {
  api_token = var.cloud_creds.cloudflare.api_token
}

provider "acme" {
  alias      = "letsencrypt"
  server_url = "https://acme-staging-v02.api.letsencrypt.org/directory"
}

provider "acme" {
  alias      = "zerossl"
  server_url = "https://acme.zerossl.com/v2/DV90"
}

provider "acme" {
  alias      = "gcp"
  server_url = "https://dv.acme-v02.api.pki.goog/directory"
}
