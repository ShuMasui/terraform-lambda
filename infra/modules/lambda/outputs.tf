output "rest_api_id" {
  description = "APIGatewayのRESTAPIリソースのIDです"
  value       = aws_api_gateway_rest_api.api.id
}