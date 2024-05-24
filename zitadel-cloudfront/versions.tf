terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.43.0"

      configuration_aliases = [aws.aws_useast]
    }

    null = {
      source  = "hashicorp/null"
      version = "3.2.2"
    }

    zitadel = {
      source  = "zitadel/zitadel"
      version = "1.1.1"
    }
  }
}
