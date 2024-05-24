variable "aws_config_file" {
  description = "The path to the AWS config file"
  type        = string
  default     = "~/.aws/config"
}

variable "aws_creds_file" {
  description = "The path to the AWS credentials file"
  type        = string
  default     = "~/.aws/credentials"
}

variable "aws_profile" {
  description = "The AWS profile to use"
  type        = string
  default     = "default"
}

variable "common_tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default = {
    Environment = "Production"
  }
}

variable "oidc_required_role" {
  description = "The name of the Zitadel role group required to access the website"
  type        = string
  default     = "authenticated"
}

variable "s3_bucket_name" {
  description = "The name of the S3 bucket to store the website content"
  type        = string
}

variable "fqdn" {
  description = "The fully qualified domain name of the website without the protocol (eg. www.example.com)"
  type        = string
}

variable "oidc_issuer" {
  description = "The OIDC issuer URL including the protocol (eg. https://login.example.com)"
  type        = string
}

variable "route53_zone_id" {
  description = "The Route53 zone ID to create the DNS record in"
  type        = string
}
