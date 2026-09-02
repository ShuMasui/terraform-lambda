terraform {
  backend "s3" {
    bucket = "study-aws-tfstate-606030504329"
    key    = "state/service"
    region = "us-east-1"
  }
}