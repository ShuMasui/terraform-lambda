variable "aws_region" {
  type        = string
  description = "AWSのリソースを配置するリージョン"
  default     = "us-east-1"
}

variable "aws_profile" {
  type        = string
  description = "AWSのプロファイル"
  default     = "study-aws"
}

variable "aws_proj" {
  type        = string
  description = "プロジェクト名"
  default     = "study-aws"
}