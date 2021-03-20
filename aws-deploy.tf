terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {
  profile = "default"
  region  = "us-east-1"
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name           = "todo"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "ID"

  attribute {
    name = "ID"
    type = "S"
  }
}

resource "aws_iam_policy" "todo-policy" {
  name        = "todo-policy"
  description = "Todo API policy to operate on AWS"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:PutItem",
          "dynamodb:DeleteItem",
          "dynamodb:GetItem",
          "dynamodb:Scan",
        ]
        Resource = [
          aws_dynamodb_table.basic-dynamodb-table.arn
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "logs:PutLogEvents",
          "logs:CreateLogGroup"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role" "todo-role" {
  name = "todo-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
  ] })
}

resource "aws_iam_role_policy_attachment" "attach-policy" {
  role       = aws_iam_role.todo-role.name
  policy_arn = aws_iam_policy.todo-policy.arn
}

resource "aws_lambda_function" "todo-lambda" {
  filename         = "function.zip"
  function_name    = "todo-lambda"
  role             = aws_iam_role.todo-role.arn
  handler          = "go-todo-lambda"
  source_code_hash = filebase64sha256("function.zip")

  runtime = "go1.x"
}

resource "aws_api_gateway_rest_api" "todo-api" {
  body = jsonencode({
    "openapi" : "3.0.1",
    "info" : {
      "title" : "todo-api",
      "description" : "Created by AWS Lambda",
      "version" : "2021-03-19T19:49:20Z"
    },
    "paths" : {
      "/todo-api" : {
        "get" : {
          "x-amazon-apigateway-integration" : {
            "httpMethod" : "POST",
            "uri" : "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/${aws_lambda_function.todo-lambda.arn}/invocations",
            "responses" : {
              "default" : {
                "statusCode" : "200"
              }
            },
            "passthroughBehavior" : "when_no_match",
            "contentHandling" : "CONVERT_TO_TEXT",
            "type" : "aws_proxy"
          }
        },
        "post" : {
          "x-amazon-apigateway-integration" : {
            "httpMethod" : "POST",
            "uri" : "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/${aws_lambda_function.todo-lambda.arn}/invocations",
            "responses" : {
              "default" : {
                "statusCode" : "200"
              }
            },
            "passthroughBehavior" : "when_no_match",
            "contentHandling" : "CONVERT_TO_TEXT",
            "type" : "aws_proxy"
          }
        }
      },
      "/todo-api/{id}" : {
        "get" : {
          "parameters" : [
            {
              "name" : "id",
              "in" : "path",
              "required" : true,
              "schema" : {
                "type" : "string"
              }
            }
          ],
          "x-amazon-apigateway-integration" : {
            "httpMethod" : "POST",
            "uri" : "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/${aws_lambda_function.todo-lambda.arn}/invocations",
            "responses" : {
              "default" : {
                "statusCode" : "200"
              }
            },
            "passthroughBehavior" : "when_no_match",
            "contentHandling" : "CONVERT_TO_TEXT",
            "type" : "aws_proxy"
          }
        },
        "delete" : {
          "parameters" : [
            {
              "name" : "id",
              "in" : "path",
              "required" : true,
              "schema" : {
                "type" : "string"
              }
            }
          ],
          "x-amazon-apigateway-integration" : {
            "httpMethod" : "POST",
            "uri" : "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/${aws_lambda_function.todo-lambda.arn}/invocations",
            "responses" : {
              "default" : {
                "statusCode" : "200"
              }
            },
            "passthroughBehavior" : "when_no_match",
            "contentHandling" : "CONVERT_TO_TEXT",
            "type" : "aws_proxy"
          }
        }
      }
    }
  })

  name = "todo-api"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "todo-deployment" {
  rest_api_id = aws_api_gateway_rest_api.todo-api.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.todo-api.body))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "todo-stage" {
  deployment_id = aws_api_gateway_deployment.todo-deployment.id
  rest_api_id   = aws_api_gateway_rest_api.todo-api.id
  stage_name    = "prod"
}

resource "aws_lambda_permission" "todo-gw-lambda" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.todo-lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.todo-api.id}/*/*/todo-api"
}

resource "aws_lambda_permission" "todo-gw-lambda-get" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.todo-lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.todo-api.id}/*/GET/*"
}

resource "aws_lambda_permission" "todo-gw-lambda-delete" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.todo-lambda.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.todo-api.id}/*/DELETE/*"
}