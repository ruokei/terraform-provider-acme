terraform {
  backend "s3" {
    region         = "ap-southeast-1"
    bucket         = "sige-test-terragrunt-s3-backend"
    key            = "multi-clouds-acme-client/basic/terraform.tfstate"
    encrypt        = true
    dynamodb_table = "sige-test-terragrunt-s3-backend"
  }
}
