resource "aws_route53_record" "frontend_svc" {
  zone_id = data.aws_route53_zone.domain.zone_id
  name    = "*.svc"
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.frontend_spa.domain_name
    zone_id                = aws_cloudfront_distribution.frontend_spa.hosted_zone_id
    evaluate_target_health = false
  }
}
