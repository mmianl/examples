output "website_cicd_aws_iam_access_id" {
  value = aws_iam_access_key.cicd.id
}

output "website_cicd_aws_iam_access_secret" {
  value = aws_iam_access_key.cicd.secret
}
