resource "aws_iam_role" "iam_for_lambda" {
  name                 = "iam_role_for_lambda"
  description          = "Allows Lambda to write to CloudWatch Logs"
  max_session_duration = 60 * 60

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  inline_policy {
    name = "writeCloudWatchLogs"
    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Sid    = "writeCloudWatchLogs"
          Effect = "Allow"
          Action = [
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ]
          Resource = "${aws_cloudwatch_log_group.this.arn}:log-stream:*"
        },
        {
          Sid    = "readSQS"
          Effect = "Allow"
          Action = [
            "sqs:ReceiveMessage",
            "sqs:DeleteMessage",
            "sqs:GetQueueAttributes"
          ]
          Resource = aws_sqs_queue.this.arn
        }
      ]
    })
  }
}