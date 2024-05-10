// ACM
resource "aws_acm_certificate" "website" {
  provider = aws.aws_useast

  domain_name       = var.fqdn
  validation_method = "DNS"
  key_algorithm     = "EC_prime256v1"

  tags = var.common_tags
}

resource "aws_route53_record" "website" {
  for_each = {
    for dvo in aws_acm_certificate.website.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = var.route53_zone_id
}

// CloudFront
resource "aws_cloudfront_distribution" "website" {
  aliases = [var.fqdn]

  origin {
    domain_name              = aws_s3_bucket.s3.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.oac.id
    origin_id                = aws_s3_bucket.s3.id
  }

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"

  web_acl_id = aws_wafv2_web_acl.docs.arn

  default_cache_behavior {
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = aws_s3_bucket.s3.id

    forwarded_values {
      query_string = true

      cookies {
        forward = "all"
      }
    }

    lambda_function_association {
      event_type   = "viewer-request"
      include_body = true
      lambda_arn   = aws_lambda_function.auth.qualified_arn
    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  price_class = "PriceClass_100"

  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = ["AT", "DE", "CH"]
    }
  }

  viewer_certificate {
    minimum_protocol_version = "TLSv1.2_2021"
    acm_certificate_arn      = aws_acm_certificate.guide.arn
    ssl_support_method       = "sni-only"
  }

  tags = var.common_tags

  depends_on = [aws_wafv2_web_acl.docs]
}

resource "aws_cloudfront_origin_access_control" "oac" {
  name                              = var.fqdn
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_route53_record" "website" {
  zone_id = var.route53_zone_id
  name    = var.fqdn
  type    = "CNAME"
  ttl     = 3600
  records = [aws_cloudfront_distribution.website.domain_name]
}
