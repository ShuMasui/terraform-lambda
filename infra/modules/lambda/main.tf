# 1. 実行者としてのlambdaに与えるIAMポリシーとそれを付与したIAMロールの作成

data "aws_iam_policy_document" "api_policy" {

  statement {
    effect = "Allow"
    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
    actions   = ["execute-api:Invoke"]
    resources = ["arn:aws:execute-api:${var.aws_region}:*:*/*/*"]
  }
}

data "aws_iam_policy_document" "lambda" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      "arn:aws:logs:*:*:*"
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "dynamodb:*"
    ]
    resources = [aws_dynamodb_table.lambda_api.arn]
  }
}

resource "aws_iam_policy" "lambda" {
  name   = var.common_name
  policy = data.aws_iam_policy_document.lambda.json
}

resource "aws_iam_role" "lambda" {
  name = var.common_name
  assume_role_policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Action : "sts:AssumeRole",
        Principal : {
          "Service" : "lambda.amazonaws.com"
        },
        Effect : "Allow",
        Sid : ""
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "api" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda.arn
}


# 2. api gatewayを作成

resource "aws_api_gateway_rest_api" "api" {
  name        = var.common_name
  description = "example serverless api"
  policy      = data.aws_iam_policy_document.api_policy.json
}

resource "aws_api_gateway_resource" "api" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "{proxy+}"
}

## ステージング先
resource "aws_api_gateway_stage" "api" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  deployment_id = aws_api_gateway_deployment.api.id
  stage_name    = "api"
}

# 3. rest api エンドポイントを定義

resource "aws_api_gateway_method" "api_get" {
  authorization = "NONE"
  http_method   = "GET"
  resource_id   = aws_api_gateway_resource.api.id
  rest_api_id   = aws_api_gateway_rest_api.api.id
}

resource "aws_api_gateway_method" "api_post" {
  authorization = "NONE"
  http_method   = "POST"
  resource_id   = aws_api_gateway_resource.api.id
  rest_api_id   = aws_api_gateway_rest_api.api.id
}

resource "aws_api_gateway_integration" "api_get" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.api.id
  http_method = aws_api_gateway_method.api_get.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.api.invoke_arn
}

resource "aws_api_gateway_integration" "api_post" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  resource_id = aws_api_gateway_resource.api.id
  http_method = aws_api_gateway_method.api_post.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.api.invoke_arn
}

resource "aws_api_gateway_deployment" "api" {
  depends_on  = [aws_api_gateway_integration.api_get, aws_api_gateway_integration.api_post]
  rest_api_id = aws_api_gateway_rest_api.api.id
}

# 4. lambda関数本体のリソースを創る

resource "aws_lambda_function" "api" {
  filename      = "functions.zip"
  function_name = "api"
  role          = aws_iam_role.lambda.arn
  runtime       = "provided.al2023"
  handler       = "bootstrap"
  memory_size   = 128
  timeout       = 900

  environment {
    variables = {
      "DYNAMO_TABLE_USERS" : aws_dynamodb_table.lambda_api.name
    }
  }
}

resource "aws_lambda_permission" "apigateway_lambda" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.api.execution_arn}/*/*/*"
}

# 5. datestoreとしてdynamodbを利用するのでそのためのリソース定義
resource "aws_dynamodb_table" "lambda_api" {
  name         = "example-table"
  billing_mode = "PAY_PER_REQUEST"

  # PKはなにか
  hash_key = "user_id"

  # テーブル定義
  attribute {
    name = "user_id"
    type = "S"
  }
}