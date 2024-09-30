resource "aws_s3_bucket" "c3_bucket" {
  bucket = "c3-eventdriven-example-bucket"
}



resource "aws_s3_bucket_policy" "allow_access_from_another_account" {
  bucket = aws_s3_bucket.c3_bucket.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::840414111995:root"
        }
        Action = [
          "s3:GetObject",
          "s3:ListBucket"
        ]
        Resource = [
          "${aws_s3_bucket.c3_bucket.arn}/*",
          aws_s3_bucket.c3_bucket.arn
        ]
      }
    ]
  })
}


resource "aws_s3_bucket_lifecycle_configuration" "this" {
  bucket = aws_s3_bucket.c3_bucket.id

  rule {
    id = "example"

    expiration {
      days                         = 90
      expired_object_delete_marker = false
    }

    filter {
      prefix = "path"
    }

    status = "Enabled"
  }
}


resource "aws_s3_object" "text" {
  bucket       = aws_s3_bucket.c3_bucket.bucket
  key          = "example.txt"
  source       = "files/example.txt"
  content_type = "text/plain"
  etag         = filemd5("files/example.txt")
}

resource "aws_s3_object" "image" {
  bucket       = aws_s3_bucket.c3_bucket.bucket
  key          = "path/icon.png"
  source       = "files/icon.png"
  content_type = "image/png"
  etag         = filemd5("files/icon.png")
}
