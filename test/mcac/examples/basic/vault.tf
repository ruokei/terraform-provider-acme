data "vault_generic_secret" "alicloud_dev" {
  path = "devops/common/cloud_credential/alicloud/dev"
}

data "vault_generic_secret" "aws_dev" {
  path = "devops/common/cloud_credential/aws/dev"
}

data "vault_generic_secret" "cloudflare_dev" {
  path = "devops/common/cloud_credential/cloudflare/dev"
}

data "vault_generic_secret" "google_dev" {
  path = "devops/common/cloud_credential/google/sige/dev"
}
