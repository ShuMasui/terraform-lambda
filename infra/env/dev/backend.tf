terraform {
  backend "s3" {
    bucket = "study-aws-tfstate-606030504329"
    key    = "state/dev"
    region = "us-east-1"
  }
}