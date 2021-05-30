resource "aws_dynamodb_table" "shortener-ddb" {
  name           = "URIStore"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "ShortID"

  attribute {
    name = "ShortID"
    type = "S"
  }
}
