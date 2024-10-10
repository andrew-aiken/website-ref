resource "aws_cloudwatch_event_rule" "console" {
  name        = "example"
  description = "Trigger lambda based on a cron expression"

  schedule_expression = "rate(10 minutes)"
}

resource "aws_cloudwatch_event_target" "sns" {
  rule      = aws_cloudwatch_event_rule.console.name
  target_id = "SendToSNS"
  arn       = aws_sns_topic.this.arn
}
