data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

data "aws_kms_key" "sns_alias" {
  key_id = "alias/aws/sns"
}

data "aws_autoscaling_group" "ecs_asg" {
  name = var.autoscaling_group_name
}

data "aws_ecs_cluster" "ecs_cluster" {
  cluster_name = var.ecs_cluster_name
}

data "archive_file" "drain_ecs_instance" {
  type        = "zip"
  source_file = "${path.module}/files/drain_ecs_instance.py"
  output_path = "${path.module}/files/drain_ecs_instance.zip"
}

data "archive_file" "complete_ecs_lifecycle" {
  type        = "zip"
  source_file = "${path.module}/files/complete_ecs_lifecycle.py"
  output_path = "${path.module}/files/complete_ecs_lifecycle.zip"
}
