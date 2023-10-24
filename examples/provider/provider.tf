terraform {
  required_providers {
    st-alicloud = {
      source = "myklst/acme"
    }
  }
}

provider "st-alicloud" {
  server_url = "https://xxx"
}
