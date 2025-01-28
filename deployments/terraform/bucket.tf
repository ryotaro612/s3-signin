// aws s3 presign s3://amzn-s3-demo-bucket/mydoc.txt --expires-in 604800 --region af-south-1 --endpoint-url https://s3.af-south-1.amazonaws.com

variable "aws_region" {
  type = string
}
variable "aws_profile" {
  type = string
}
variable "bucket_name" {
  type = string
}
variable "key" {
  type = string
}


provider "aws" {
  region = var.aws_region
  profile = var.aws_profile
}

resource "aws_s3_bucket" "sample_bucket" {
  bucket = var.bucket_name

}

resource "aws_s3_object" "sample_object" {
  bucket = var.bucket_name
  key    = var.key
  source = "index.txt"
}

