resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/${local.function_name}"
  retention_in_days = 30
  skip_destroy      = false
}

# trivy:ignore:avd-aws-0066
resource "aws_lambda_function" "this" {
  function_name = local.function_name
  description   = "Example description"
  runtime       = "nodejs18.x"
  package_type  = "Zip"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "index.handler"
  memory_size   = 128
  timeout       = 30

  source_code_hash = data.archive_file.this.output_base64sha256
  filename         = data.archive_file.this.output_path

  environment {
    variables = {
      defaultMessage = "hello world"
    }
  }
}
