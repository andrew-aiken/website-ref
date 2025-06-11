resource "aws_acm_certificate" "svc" {
  lifecycle {
    create_before_destroy = true
  }

  domain_name       = "*.svc.${local.domain}"
  validation_method = "DNS"

  subject_alternative_names = [
    "*.svc.${local.domain}"
  ]
}

resource "aws_acm_certificate_validation" "svc_validation" {
  certificate_arn = aws_acm_certificate.svc.arn
  validation_record_fqdns = [
    for record in aws_acm_certificate.svc.domain_validation_options : record.resource_record_name
  ]
}

resource "aws_route53_record" "svc_validation" {
  for_each = {
    for idx, record in aws_acm_certificate.svc.domain_validation_options :
    record.domain_name => record
  }

  allow_overwrite = true
  name            = each.value.resource_record_name
  records         = [each.value.resource_record_value]
  ttl             = 300
  type            = each.value.resource_record_type
  zone_id         = data.aws_route53_zone.domain.zone_id
}
