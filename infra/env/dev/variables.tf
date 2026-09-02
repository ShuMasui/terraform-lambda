variable "aws_region" {
  type        = string
  description = "AWSのリソースを配置するリージョン"
  default     = "us-east-1"
}

variable "aws_profile" {
  type        = string
  description = "AWSのプロファイル"
}

variable "aws_proj" {
  type        = string
  description = "プロジェクト名"
}

variable "vpc_cidr" {
  type        = string
  description = "VPCにおけるCIDR表記されたネットワーク範囲"
}

variable "env" {
  type        = string
  description = "環境"
}