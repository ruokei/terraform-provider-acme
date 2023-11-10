module "multi_clouds_acme_client" {
  source = "../../"

  cloud_creds = {
    alicloud = {
      access_key = data.vault_generic_secret.alicloud_dev.data["ALICLOUD_ACCESS_KEY"]
      secret_key = data.vault_generic_secret.alicloud_dev.data["ALICLOUD_SECRET_KEY"]
    }
    aws = {
      access_key = data.vault_generic_secret.aws_dev.data["AWS_ACCESS_KEY_ID"]
      secret_key = data.vault_generic_secret.aws_dev.data["AWS_SECRET_ACCESS_KEY"]
    }
    cloudflare = {
      api_token = data.vault_generic_secret.cloudflare_dev.data["CLOUDFLARE_API_TOKEN"]
    }
    google = {
      project     = data.vault_generic_secret.google_dev.data["GOOGLE_PROJECT"]
      credentials = data.vault_generic_secret.google_dev.data["GOOGLE_APPLICATION_CREDENTIALS"]
    }
  }

  module_info = {
    brand                  = "sige"
    env                    = "basic"
    be_app_category        = "generic"
    alicloud_region = "ap-southeast-1"
    aws_region             = "ap-southeast-1"
    gcp_region             = "asia-east2"
  }
  module_tmpl = {}

  alicloud_oss_regions = ["cn-hongkong", "eu-west-1"]
  aws_s3_regions       = ["ap-southeast-1"]
  gcp_storage_regions  = ["ASIA", "EU", "US"]

  acme_client = {
    renew_before_days      = 90
    register_email_address = "test@sige.la"
    ca_list                = ["letsencrypt", "zerossl", "gcp"]

    domain_sans = {
      "sige-test.com"  = [],
      "sige-test2.com" = []
    }
  }

  extra_domains = var.extra_domains
}
