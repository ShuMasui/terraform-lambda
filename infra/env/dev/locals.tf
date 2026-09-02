locals {
  common_tags = {
    Name    = "${var.env}-${var.aws_proj}"
    Env     = var.env
    Project = var.aws_proj
  }

  common_tags_public = {
    Name    = "${var.env}-${var.aws_proj}-public"
    Env     = var.env
    Project = var.aws_proj
  }

  common_tags_private = {
    Name    = "${var.env}-${var.aws_proj}-private"
    Env     = var.env
    Project = var.aws_proj
  }
}