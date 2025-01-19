variable "aws_region" {
  type = string
}
variable "aws_profile" {
  type = string
}
variable "bucket_name" {
  type = string
}


provider "aws" {
  region = var.aws_region
  profile = var.aws_profile
}

resource "aws_s3_bucket" "example" {
  bucket = var.bucket_name

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}
