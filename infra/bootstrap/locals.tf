locals {
  state_bucket_name = "${var.aws_proj}-tfstate-${data.aws_caller_identity.current.account_id}"
}