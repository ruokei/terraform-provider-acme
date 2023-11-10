resource "alicloud_oss_bucket" "oss_bucket" {
  bucket        = "${var.resource_name}-${var.ca}-${var.region}"
  acl           = "private"
  force_destroy = true

  tags = {
    app = var.app_tag
    ca  = var.ca
  }

  server_side_encryption_rule {
    sse_algorithm = "AES256"
  }
}

##################################################################################

resource "alicloud_ram_policy" "ram_policy" {
  policy_name     = "${var.resource_name}-${var.ca}-${var.region}"
  rotate_strategy = "DeleteOldestNonDefaultVersionWhenLimitExceeded"
  force           = true
  policy_document = <<EOF
{
  "Version": "1",
  "Statement": [
    {
      "Action": [
        "oss:ListBuckets",
        "oss:GetBucketStat",
        "oss:GetBucketInfo",
        "oss:GetBucketAcl",
        "oss:GetBucketTagging"
      ],
      "Effect": "Allow",
      "Resource": [
        "acs:oss:*:*:*"
      ]
    },
    {
      "Action": [
        "oss:HeadObject",
        "oss:GetObject",
        "oss:PutObject",
        "oss:ListObjects",
        "oss:GetObjectTagging",
        "oss:PutObjectTagging"
      ],
      "Effect": "Allow",
      "Resource": [
        "acs:oss:*:*:${alicloud_oss_bucket.oss_bucket.id}",
        "acs:oss:*:*:${alicloud_oss_bucket.oss_bucket.id}/*"
      ]
    }
  ]
}
EOF
}
