# trivy:ignore:avd-aws-0090
resource "aws_s3_bucket" "frontend_spa" {
  bucket = "example-spa-frontend"
}

resource "aws_s3_bucket_public_access_block" "frontend_spa" {
  bucket = aws_s3_bucket.frontend_spa.id

  block_public_acls       = true
  ignore_public_acls      = true
  block_public_policy     = true
  restrict_public_buckets = true
}

# trivy:ignore:avd-aws-0132
resource "aws_s3_bucket_server_side_encryption_configuration" "frontend_spa_sse_s3_encryption" {
  bucket = aws_s3_bucket.frontend_spa.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
    bucket_key_enabled = true
  }
}

resource "aws_s3_bucket_policy" "frontend_spa" {
  bucket = aws_s3_bucket.frontend_spa.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "RequireEncryptedConnections"
        Action = "s3:*"
        Effect = "Deny"
        Resource = [
          aws_s3_bucket.frontend_spa.arn,
          "${aws_s3_bucket.frontend_spa.arn}/*"
        ]
        Condition = {
          Bool = {
            "aws:SecureTransport" = false
          }
        }
        Principal = "*"
      },
      {
        Sid    = "AllowCloudFrontServicePrincipal"
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action   = "s3:GetObject"
        Resource = "${aws_s3_bucket.frontend_spa.arn}/main/*"
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.frontend_spa.arn
          }
        }
      }
    ]
  })
}

resource "aws_s3_bucket_lifecycle_configuration" "frontend_spa" {
  provider = aws.primary_region

  bucket = aws_s3_bucket.frontend_spa.id

  rule {
    id     = "ExpireOldVersions"
    status = "Enabled"


    expiration {
      days                         = 1
      expired_object_delete_marker = false
    }

    filter {
      and {
        prefix = "main/"
        tags = {
          "expire" = "true"
        }
      }
    }
  }
}
