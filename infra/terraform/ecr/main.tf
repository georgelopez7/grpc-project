provider "aws" {
  region = "eu-west-2"
}

resource "aws_ecr_repository" "services" {
  for_each = toset(var.service_names)

  name = each.key

  image_scanning_configuration {
    scan_on_push = true
  }

  image_tag_mutability = "MUTABLE"

  encryption_configuration {
    encryption_type = "AES256"
  }
}