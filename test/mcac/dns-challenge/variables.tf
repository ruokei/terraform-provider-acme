variable "domain" {
  description = "Domain to run checking"
  type        = string
}

variable "cloud_creds" {
  description = <<EOL
  Cloud credentials.
    - `alicloud`       : AliCloud credential.
      - `access_key`   : AliCloud non-CN access key.
      - `secret_key`   : AliCloud non-CN secret key.
    - `aws`            : AWS credential.
      - `region`       : AWS region.
      - `access_key`   : AWS access key.
      - `secret_key`   : AWS secret key.
    - `cloudflare`     : Cloudflare credential.
      - `api_token`    : Cloudflare API Token.
  EOL
  type = object({
    alicloud = object({
      access_key = string
      secret_key = string
    })
    aws = object({
      region     = string
      access_key = string
      secret_key = string
    })
    cloudflare = object({
      api_token = string
    })
  })
}
