resource "aws_sns_topic" "this" {
  name = "example"
}

resource "aws_sns_topic_subscription" "trigger_lambda" {
  topic_arn = aws_sns_topic.this.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.this.arn
}
