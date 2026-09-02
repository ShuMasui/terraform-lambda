resource "aws_vpc" "vpc" {
  cidr_block = var.vpc_cidr

  tags = var.common_tags

}

resource "aws_subnet" "public_subnet" {

  for_each = local.subnet_numbers
  vpc_id   = aws_vpc.vpc.id

  availability_zone = each.key
  cidr_block        = cidrsubnet(aws_vpc.vpc.cidr_block, 8, each.value)

  tags = var.common_tags_public

}

resource "aws_subnet" "private_subnet" {

  for_each = local.subnet_numbers
  vpc_id   = aws_vpc.vpc.id

  availability_zone = each.key
  cidr_block        = cidrsubnet(aws_vpc.vpc.cidr_block, 8, each.value + 3)

  tags = var.common_tags_private

}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id

  tags = var.common_tags
}

resource "aws_eip" "nat_gateway" {
  depends_on = [aws_internet_gateway.igw]

  tags = var.common_tags
}

resource "aws_nat_gateway" "nat_gateway" {
  allocation_id = aws_eip.nat_gateway.id
  subnet_id     = aws_subnet.public_subnet["us-east-1a"].id
  depends_on    = [aws_internet_gateway.igw]

  tags = var.common_tags
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.vpc.id

  tags = var.common_tags_public
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.vpc.id

  tags = var.common_tags_private
}

resource "aws_route_table_association" "public" {
  for_each       = local.subnet_numbers
  subnet_id      = aws_subnet.public_subnet[each.key].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  for_each       = local.subnet_numbers
  subnet_id      = aws_subnet.private_subnet[each.key].id
  route_table_id = aws_route_table.private.id
}

resource "aws_route" "public" {
  route_table_id         = aws_route_table.public.id
  gateway_id             = aws_internet_gateway.igw.id
  destination_cidr_block = "0.0.0.0/0"
}

resource "aws_route" "private" {
  route_table_id         = aws_route_table.private.id
  gateway_id             = aws_nat_gateway.nat_gateway.id
  destination_cidr_block = "0.0.0.0/0"
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id       = aws_vpc.vpc.id
  service_name = "com.amazonaws.${var.aws_region}.s3"
  tags         = var.common_tags
}

resource "aws_vpc_endpoint_route_table_association" "private_s3" {
  vpc_endpoint_id = aws_vpc_endpoint.s3.id
  route_table_id  = aws_route_table.private.id
}