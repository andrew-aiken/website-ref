resource "aws_lambda_function" "drain_ecs_instance" {
  function_name    = var.drain_ecs_instance_name
  description      = "Drain tasks from ECS instances"
  handler          = "drain_ecs_instance.lambda_handler"
  role             = aws_iam_role.execution_drain_ecs_instance.arn
  runtime          = "python3.11"
  memory_size      = 128
  timeout          = 3
  filename         = data.archive_file.drain_ecs_instance.output_path
  source_code_hash = data.archive_file.drain_ecs_instance.output_base64sha256
  environment {
    variables = {
      ClusterName = var.ecs_cluster_name
      LOGLEVEL    = "info"
    }
  }
  tags = var.tags
}

resource "aws_cloudwatch_log_group" "lambda_log_group" {
  name              = "/aws/lambda/${var.drain_ecs_instance_name}"
  retention_in_days = 0

  lifecycle {
    prevent_destroy = false
  }

  tags = var.tags
}

resource "aws_lambda_permission" "drain_ecs_instance_function" {
  function_name = aws_lambda_function.drain_ecs_instance.arn
  statement_id  = "LambdaInvokePermission"
  action        = "lambda:InvokeFunction"
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.drain_ecs_instance.arn
}


resource "aws_iam_role" "execution_drain_ecs_instance" {
  name        = "${var.ecs_cluster_name}-drain-ecs-instance-lambda-execution-role"
  description = "Allows ECS tasks to remove docker containers before terminate"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = "arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.id}:function:${var.drain_ecs_instance_name}"
          }
        },
        Effect = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy" "execution_drain_ecs_instance_policy" {
  name = "${var.ecs_cluster_name}-drain-ecs-instance-sns-lambda-policy"
  role = aws_iam_role.execution_drain_ecs_instance.name
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Effect = "Allow",
        Resource = [
          "${aws_cloudwatch_log_group.lambda_log_group.arn}:*"
        ],
        Sid = "Log"
      },
      {
        Action = [
          "ecs:ListContainerInstances",
          "ecs:DescribeContainerInstances",
          "ecs:UpdateContainerInstancesState",
          "ecs:TagResource"
        ],
        Effect = "Allow",
        Resource = [
          data.aws_ecs_cluster.ecs_cluster.arn,
          "arn:aws:ecs:${data.aws_region.current.name}:${data.aws_caller_identity.current.id}:container-instance/${var.ecs_cluster_name}/*"
        ],
        Sid = "asg"
      }
    ]
  })
}

resource "aws_sns_topic" "drain_ecs_instance" {
  name              = "${var.ecs_cluster_name}-drain-ecs-instance"
  kms_master_key_id = data.aws_kms_key.sns_alias.key_id


  tags = var.tags
}

resource "aws_sns_topic_subscription" "drain_ecs_instance_sns" {
  topic_arn = aws_sns_topic.drain_ecs_instance.arn
  protocol  = "lambda"
  endpoint  = aws_lambda_function.drain_ecs_instance.arn
}

resource "aws_autoscaling_lifecycle_hook" "drain_ecs_instance_hook" {
  name                    = "${var.ecs_cluster_name}-drain-ecs-instance"
  autoscaling_group_name  = data.aws_autoscaling_group.ecs_asg.id
  default_result          = "ABANDON"
  heartbeat_timeout       = 900
  lifecycle_transition    = "autoscaling:EC2_INSTANCE_TERMINATING"
  notification_target_arn = aws_sns_topic.drain_ecs_instance.arn
  role_arn                = aws_iam_role.sns_drain_ecs_instance.arn
}

resource "aws_iam_role" "sns_drain_ecs_instance" {
  name        = "${var.ecs_cluster_name}-drain-ecs-instance-sns"
  description = "Allows ECS tasks to remove docker containers before terminate"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "autoscaling.amazonaws.com"
        }
      }
    ]
  })

  inline_policy {
    name = "sns-public"

    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action   = "sns:Publish"
          Effect   = "Allow"
          Resource = aws_sns_topic.drain_ecs_instance.arn
        }
      ]
    })
  }

  tags = var.tags
}
