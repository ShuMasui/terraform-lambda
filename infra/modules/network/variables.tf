variable "vpc_cidr" {
  type        = string
  description = "VPCのサブネット分割を含めたCIDR表記でのネットワーク範囲を決めるものです"
}

variable "common_tags" {
  type        = map(string)
  description = "標準的なタグ表記"
}

variable "common_tags_public" {
  type        = map(string)
  description = "標準的なタグ表記（公開領域）"
}

variable "common_tags_private" {
  type        = map(string)
  description = "標準的なタグ表記（非公開領域）"
}

variable "aws_region" {
  type        = string
  description = "AWSリソースの設置地域"
}