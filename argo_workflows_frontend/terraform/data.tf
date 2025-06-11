data "aws_route53_zone" "domain" {
  name         = local.domain
  private_zone = false
}


data "aws_cloudfront_cache_policy" "cache_optimized" {
  name = "Managed-CachingOptimized"
}
