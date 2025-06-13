resource "aws_cloudfront_function" "frontend_spa" {
  name    = "frontend-spa-domain-rewrite"
  runtime = "cloudfront-js-2.0"
  comment = "rewrites domains to s3 paths and supports SPA react routing"
  publish = true
  code    = file("${path.module}/src/frontend.js")
}


resource "aws_cloudfront_response_headers_policy" "cors_policy" {
  name    = "frontend-spa-cors-policy"
  comment = "CORS policy for frontend SPA"

  cors_config {
    access_control_allow_credentials = false
    access_control_allow_headers {
      items = ["*"]
    }
    access_control_allow_methods {
      items = ["GET", "HEAD", "OPTIONS"]
    }
    access_control_allow_origins {
      items = ["https://app.${local.domain}"]
    }
    access_control_expose_headers {
      items = ["*"]
    }
    access_control_max_age_sec = 600
    origin_override            = true
  }
}


resource "aws_cloudfront_distribution" "frontend_spa" {
  comment = "Serves frontend service SPAs"

  aliases = [
    "*.svc.${local.domain}"
  ]

  origin {
    connection_attempts      = 3
    connection_timeout       = 10
    domain_name              = aws_s3_bucket.frontend_spa.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.s3_origin.id
    origin_id                = aws_s3_bucket.frontend_spa.bucket_regional_domain_name
    origin_path              = "/main"
  }

  default_cache_behavior {
    allowed_methods = [
      "GET",
      "HEAD",
      "OPTIONS"
    ]
    cached_methods = [
      "GET",
      "HEAD",
      "OPTIONS"
    ]
    target_origin_id           = aws_s3_bucket.frontend_spa.bucket_regional_domain_name
    cache_policy_id            = data.aws_cloudfront_cache_policy.cache_optimized.id
    compress                   = true
    viewer_protocol_policy     = "redirect-to-https"
    response_headers_policy_id = aws_cloudfront_response_headers_policy.cors_policy.id

    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.frontend_spa.arn
    }
  }

  price_class     = "PriceClass_100"
  enabled         = true
  http_version    = "http2"
  is_ipv6_enabled = true

  viewer_certificate {
    acm_certificate_arn            = aws_acm_certificate.svc.arn
    cloudfront_default_certificate = false
    minimum_protocol_version       = "TLSv1.2_2021"
    ssl_support_method             = "sni-only"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }
}

