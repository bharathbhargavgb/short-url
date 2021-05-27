resource "aws_api_gateway_rest_api" "shortener-api-gateway" {
  name        = "shortener-api-gateway"
  description = "API gateway for REST APIs to URL shortener service"
}

resource "aws_api_gateway_resource" "expand-resource" {
  rest_api_id = aws_api_gateway_rest_api.shortener-api-gateway.id
  parent_id   = aws_api_gateway_rest_api.shortener-api-gateway.root_resource_id
  path_part   = "expand"
}

resource "aws_api_gateway_resource" "shorten-resource" {
  rest_api_id = aws_api_gateway_rest_api.shortener-api-gateway.id
  parent_id   = aws_api_gateway_rest_api.shortener-api-gateway.root_resource_id
  path_part   = "shorten"
}

resource "aws_api_gateway_method" "expand-method" {
  rest_api_id   = aws_api_gateway_rest_api.shortener-api-gateway.id
  resource_id   = aws_api_gateway_resource.expand-resource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_method" "shorten-method" {
  rest_api_id   = aws_api_gateway_rest_api.shortener-api-gateway.id
  resource_id   = aws_api_gateway_resource.shorten-resource.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "expand-lambda-integration" {
  rest_api_id = aws_api_gateway_rest_api.shortener-api-gateway.id
  resource_id = aws_api_gateway_resource.expand-resource.id
  http_method = aws_api_gateway_method.expand-method.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.shortener-lambda.invoke_arn
}

resource "aws_api_gateway_integration" "shorten-lambda-integration" {
  rest_api_id = aws_api_gateway_rest_api.shortener-api-gateway.id
  resource_id = aws_api_gateway_resource.shorten-resource.id
  http_method = aws_api_gateway_method.shorten-method.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.shortener-lambda.invoke_arn
}

resource "aws_lambda_permission" "lambda-permission" {
  statement_id  = "AllowAPIGatewayToInvokeURLShortenerLambda"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.shortener-lambda.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.shortener-api-gateway.execution_arn}/*/*/*"
}

resource "aws_api_gateway_deployment" "shortener-deployment" {
  depends_on = [
    aws_api_gateway_integration.expand-lambda-integration,
    aws_api_gateway_integration.shorten-lambda-integration,
  ]

  rest_api_id = aws_api_gateway_rest_api.shortener-api-gateway.id
  stage_name  = "beta"
}

output "base_url" {
  value = aws_api_gateway_deployment.shortener-deployment.invoke_url
}
