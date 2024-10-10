resource "aws_sqs_queue" "this" {
  name                      = "example-queue"
  receive_wait_time_seconds = 10
}

resource "aws_lambda_event_source_mapping" "example" {
  event_source_arn = aws_sqs_queue.this.arn
  function_name    = aws_lambda_function.this.arn
}

resource "aws_lambda_permission" "sqs" {
  statement_id  = "AllowExecutionFromSQS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "sqs.amazonaws.com"
  source_arn    = aws_sqs_queue.this.arn
}
