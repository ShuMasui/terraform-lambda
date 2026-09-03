
# module "network" {
#   source = "../../modules/network"

#   vpc_cidr = var.vpc_cidr

#   common_tags = local.common_tags

#   common_tags_public = local.common_tags_public

#   common_tags_private = local.common_tags_private

#   aws_region = var.aws_region
# }

module "lambda" {
  source = "../../modules/lambda"

  aws_region = var.aws_region

  common_name = "lambda"
}

#============================================================================
# Github Actions 用ロール
#============================================================================
# -----------------------------------------------------------------------------
# GitHub Actions プロバイダー設定
# -----------------------------------------------------------------------------
resource "aws_iam_openid_connect_provider" "terraform_cicd" {
  url             = "https://token.actions.githubusercontent.com"
  client_id_list  = ["sts.amazonaws.com"]
  # このコードは固定値
  # OIDC ID プロバイダーのサムプリント
  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"]
}

# -----------------------------------------------------------------------------
# GitHub Actions 用ロール作成
# -----------------------------------------------------------------------------
resource "aws_iam_role" "terraform_cicd_oidc_role" {
  name = "TerraCICDDemoOIDCRole"
  path = "/"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = "sts:AssumeRoleWithWebIdentity"
      Principal = {
        Federated = aws_iam_openid_connect_provider.terraform_cicd.arn
      }
      Condition = {
        StringLike = {
          "token.actions.githubusercontent.com:sub" = [
            # GitHubが八虎するトークンでアクセスできるリポジトリ
            # 指定したリポジトリの全リソースにアクセス
            "repo:ShuMasui/terraform-lambda:*",
          ]
        }
      }
    }]
  })
}

# -----------------------------------------------------------------------------
# ポリシーのアタッチ（AdministratorAccess_attachment）
# -----------------------------------------------------------------------------
resource "aws_iam_role_policy_attachment" "AdministratorAccess_attachment" {
  role       = aws_iam_role.terraform_cicd_oidc_role.name
  # Admin権限を指定
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}
