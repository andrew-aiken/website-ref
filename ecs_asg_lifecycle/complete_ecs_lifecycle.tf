resource "aws_lambda_function" "complete_ecs_lifecycle" {
  function_name    = var.complete_ecs_lifecycle_name
  description      = "Mark ECS ASG as lifecycle as completed, letting the instances to be terminated by autoscaling"
  handler          = "complete_ecs_lifecycle.lambda_handler"
  role             = aws_iam_role.complete_ecs_lifecycle.arn
  runtime          = "python3.11"
  memory_size      = 128
  timeout          = 3
  filename         = data.archive_file.complete_ecs_lifecycle.output_path
  source_code_hash = data.archive_file.complete_ecs_lifecycle.output_base64sha256

  environment {
    variables = {
      LOGLEVEL = "info"
    }
  }

  tags = var.tags
}

resource "aws_cloudwatch_log_group" "complete_ecs_lifecycle_log_group" {
  name              = "/aws/lambda/${var.complete_ecs_lifecycle_name}"
  retention_in_days = 0

  lifecycle {
    prevent_destroy = false
  }

  tags = var.tags
}

resource "aws_lambda_permission" "eventbridge_complete_ecs_lifecycle_function" {
  function_name = aws_lambda_function.complete_ecs_lifecycle.arn
  statement_id  = "LambdaInvokePermission"
  action        = "lambda:InvokeFunction"
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.complete_ecs_lifecycle.arn
}

resource "aws_iam_role" "complete_ecs_lifecycle" {
  name                 = "${var.ecs_cluster_name}-drain-ecs-instance"
  description          = "Role used to mark ECS autoscaling lifecycles as complete"
  max_session_duration = 60 * 60

  managed_policy_arns = [
    aws_iam_policy.complete_ecs_lifecycle.id
  ]

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = "arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.id}:function:${var.complete_ecs_lifecycle_name}"
          }
        }
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Sid = ""
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_policy" "complete_ecs_lifecycle" {
  name        = "${var.ecs_cluster_name}-complete-ecs-lifecycle"
  path        = "/"
  description = "Policy that allows the marking of ECS autoscaling lifecycles as complete"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Effect   = "Allow"
        Resource = "${aws_cloudwatch_log_group.complete_ecs_lifecycle_log_group.arn}:*"
        Sid      = "log"
      },
      {
        Action = [
          "ec2:DescribeInstances",
          "autoscaling:DescribeAutoScalingGroups",
          "autoscaling:DescribeLifecycleHooks"
        ]
        Effect   = "Allow"
        Resource = "*"
        Sid      = "ec2"
      },
      {
        Action = [
          "autoscaling:CompleteLifecycleAction"
        ]
        Effect   = "Allow"
        Resource = data.aws_autoscaling_group.ecs_asg.arn
        Sid      = "asg"
      }
    ]
  })

  tags = var.tags
}

resource "aws_cloudwatch_event_rule" "complete_ecs_lifecycle" {
  name        = "${var.ecs_cluster_name}-complete-ecs-lifecycle"
  description = "Triggers a lambda function when a ECS instance is draining and has no tasks running"

  is_enabled = true

  event_pattern = jsonencode({
    source      = ["aws.ecs"]
    detail-type = ["ECS Container Instance State Change"]
    detail = {
      clusterArn        = [data.aws_ecs_cluster.ecs_cluster.arn]
      pendingTasksCount = [0]
      runningTasksCount = [0]
      status            = ["DRAINING"]
    }
  })

  tags = var.tags
}

resource "aws_cloudwatch_event_target" "complete_ecs_lifecycle" {
  rule      = aws_cloudwatch_event_rule.complete_ecs_lifecycle.name
  target_id = "trigger-complete-ecs-lifecycle-lambda"
  arn       = aws_lambda_function.complete_ecs_lifecycle.arn
}
