
module "network" {
  source = "../../modules/network"

  vpc_cidr = var.vpc_cidr

  common_tags = local.common_tags

  common_tags_public = local.common_tags_public

  common_tags_private = local.common_tags_private

  aws_region = var.aws_region
}