locals {
  redirect_uri                = "https://${var.fqdn}/auth/callback"
  post_logout_redirect_uri    = "https://${var.fqdn}/"
  oidc_jwks_uri               = "${var.oidc_issuer}/oauth/v2/keys"
  oidc_token_endpoint         = "${var.oidc_issuer}/oauth/v2/token"
  oidc_authorization_endpoint = "${var.oidc_issuer}/oauth/v2/authorize"
  oidc_role_key               = "urn:zitadel:iam:org:project:roles"
}
