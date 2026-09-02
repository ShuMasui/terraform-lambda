output "state_bucket_name" {
  description = "tfstate用S3バケット名。env/*/backend.tf のbucketに設定する"
  value       = aws_s3_bucket.terraform_state.id
}
