locals {
  lambda-payload-path = "${path.module}${var.lambda-payload-relative-path}${var.lambda-zip-name}"
}

resource "aws_iam_role" "shortener-lambda-role" {
  name = "shortener-lambda-role"
  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Principal" : {
          "Service" : "lambda.amazonaws.com"
        },
        "Action" : "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_lambda_function" "shortener-lambda" {
  filename      = "${local.lambda-payload-path}"
  function_name = "url-shortener"
  role          = aws_iam_role.shortener-lambda-role.arn
  handler       = "main"

  source_code_hash = filebase64sha256("${local.lambda-payload-path}")

  runtime = "go1.x"

  environment {
    variables = {
      URI_STORE = "URIStore"
    }
  }
}

data "aws_iam_policy" "AWSLambdaBasicExecutionRole" {
  arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// TODO: Improve the structure: https://stackoverflow.com/a/45487184/3505471
resource "aws_iam_role_policy_attachment" "lambda-logs" {
  role       = aws_iam_role.shortener-lambda-role.name
  policy_arn = data.aws_iam_policy.AWSLambdaBasicExecutionRole.arn
}

resource "aws_iam_policy" "dynamodb-CRUD-role" {
  name = "dynamodb-CRUD-role"
  policy = jsonencode({
    "Statement" : [
      {
        "Action" : [
          "dynamodb:PutItem",
          "dynamodb:GetItem"
        ],
        "Effect" : "Allow",
        "Resource" : "*"
      }
    ],
    "Version" : "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "dynamodb-CRUD" {
  role       = aws_iam_role.shortener-lambda-role.name
  policy_arn = aws_iam_policy.dynamodb-CRUD-role.arn
}

