variable "drain_ecs_instance_name" {
  type        = string
  default     = "drain-ecs-instance"
  description = "Name that resource for the draining lambda function should use"
}

variable "complete_ecs_lifecycle_name" {
  type        = string
  default     = "complete-ecs-lifecycle"
  description = "Name that resource for function that completes the lifecycle should use"
}

variable "autoscaling_group_name" {
  type        = string
  description = "Name of the ECS ASG"
}

variable "ecs_cluster_name" {
  type        = string
  description = "Name of the ECS cluster"
}

variable "tags" {
  type        = map(string)
  description = "Map of tags to be applied to all resources in cluster"
  default     = {}
}
