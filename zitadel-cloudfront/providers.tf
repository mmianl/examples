provider "aws" {
  shared_config_files      = [var.aws_config_file]
  shared_credentials_files = [var.aws_creds_file]
  profile                  = var.aws_profile
}

provider "aws" {
  alias = "aws_useast"

  shared_config_files      = [var.aws_config_file]
  shared_credentials_files = [var.aws_creds_file]
  profile                  = var.aws_profile
  region                   = "us-east-1"
}

provider "null" {}

provider "zitadel" {
  domain           = data.terraform_remote_state.auth.outputs.zitadel_fqdn
  insecure         = false
  port             = "443"
  jwt_profile_file = "jwt-profile.json"
}
