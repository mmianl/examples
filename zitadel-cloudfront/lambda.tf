resource "local_file" "oidc_config" {
  filename        = "${path.module}/lambda/src/${var.s3_bucket_name}.json"
  file_permission = "0664"

  content = templatefile("${path.module}/config.json.tftpl", {
    oidc_issuer                 = var.oidc_issuer
    oidc_client_id              = zitadel_application_oidc.website.client_id
    oidc_redirect_uri           = local.redirect_uri
    oidc_jwks_uri               = local.oidc_jwks_uri
    oidc_token_endpoint         = local.oidc_token_endpoint
    oidc_authorization_endpoint = local.oidc_authorization_endpoint
    oidc_role_key               = local.oidc_role_key
    oidc_required_role          = var.oidc_required_role
  })
}

resource "null_resource" "npm_build" {
  triggers = {
    dir_sha1           = sha1(join("", [for f in fileset("", "${path.module}/lambda/src/*.ts") : filesha1(f)]))
    template           = file("${path.module}/config.json.tftpl")
    oidc_issuer        = var.oidc_issuer
    oidc_client_id     = zitadel_application_oidc.website.client_id
    oidc_redirect_uri  = local.redirect_uri
    oidc_required_role = var.oidc_required_role
  }

  provisioner "local-exec" {
    working_dir = "${path.module}/lambda"
    command     = "NAME=${var.s3_bucket_name} npm run build"
  }

  depends_on = [
    local_file.oidc_config
  ]
}

data "archive_file" "lambda" {
  type        = "zip"
  source_dir  = "${path.module}/lambda/${var.s3_bucket_name}-dist/"
  output_path = "${path.module}/lambda/${var.s3_bucket_name}-dist.zip"

  depends_on = [
    null_resource.npm_build
  ]
}

data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type = "Service"
      identifiers = [
        "lambda.amazonaws.com",
        "edgelambda.amazonaws.com"
      ]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "lambda_auth" {
  name               = "${var.fqdn}-iam-role"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_iam_policy" "lambda_logging" {
  name = "${var.fqdn}-logging-policy"
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        Action : [
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:CreateLogGroup",
        ],
        Effect : "Allow",
        Resource : "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_logging" {
  role       = aws_iam_role.lambda_auth.id
  policy_arn = aws_iam_policy.lambda_logging.arn
}

resource "aws_cloudwatch_log_group" "lambda_log_group" {
  provider = aws.aws_useast

  name              = "/aws/lambda/lambdaauthn"
  retention_in_days = 7

  tags = var.common_tags
}

resource "aws_lambda_function" "auth" {
  provider = aws.aws_useast

  filename      = "${path.module}/lambda/${var.s3_bucket_name}-dist.zip"
  function_name = "lambdaauthn"
  role          = aws_iam_role.lambda_auth.arn
  handler       = "index.handler"
  memory_size   = 128 // MB
  timeout       = 5   // Seconds
  publish       = true
  architectures = ["x86_64"]

  source_code_hash = data.archive_file.lambda.output_base64sha256

  runtime = "nodejs18.x"

  tags = var.common_tags

  depends_on = [
    data.archive_file.lambda,
  ]
}
