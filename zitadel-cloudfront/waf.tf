# WAF
resource "aws_wafv2_web_acl" "docs" {
  provider = aws.aws_useast

  name        = var.fqdn
  description = "Cloudfront rate limiting"
  scope       = "CLOUDFRONT"

  default_action {
    allow {}
  }

  rule {
    name     = "bot-control"
    priority = 1

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesBotControlRuleSet"
        vendor_name = "AWS"

        managed_rule_group_configs {

          aws_managed_rules_bot_control_rule_set {
            inspection_level = "COMMON"
          }
        }
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "${var.fqdn}-waf-bot-control"
      sampled_requests_enabled   = false
    }
  }

  rule {
    name     = "rate-limiting"
    priority = 2

    action {
      block {}
    }

    statement {
      rate_based_statement {
        limit              = 1000
        aggregate_key_type = "IP"

        scope_down_statement {
          geo_match_statement {
            country_codes = ["AT", "DE", "CH"]
          }
        }
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "${var.fqdn}-waf-rate-limiting"
      sampled_requests_enabled   = false
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name                = "${var.fqdn}-waf"
    sampled_requests_enabled   = false
  }

  tags = var.common_tags
}

resource "aws_cloudwatch_log_group" "docs" {
  provider = aws.aws_useast

  retention_in_days = 7
  name              = "aws-waf-logs-${var.fqdn}"

  tags = var.common_tags
}

resource "aws_wafv2_web_acl_logging_configuration" "docs" {
  provider = aws.aws_useast

  log_destination_configs = [aws_cloudwatch_log_group.docs.arn]
  resource_arn            = aws_wafv2_web_acl.docs.arn

  redacted_fields {
    single_header {
      name = "cookie"
    }
  }
}
