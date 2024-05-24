resource "zitadel_org" "website" {
  name = "Website"
}

resource "zitadel_project" "website" {
  name                     = "Website"
  org_id                   = zitadel_org.website.id
  project_role_assertion   = true
  project_role_check       = false
  has_project_check        = false
  private_labeling_setting = "PRIVATE_LABELING_SETTING_UNSPECIFIED"
}

resource "zitadel_application_oidc" "website" {
  project_id                  = zitadel_project.website.id
  org_id                      = zitadel_org.website.id
  name                        = "Website"
  redirect_uris               = [local.redirect_uri]
  response_types              = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types                 = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  post_logout_redirect_uris   = [local.post_logout_redirect_uri]
  app_type                    = "OIDC_APP_TYPE_WEB"
  auth_method_type            = "OIDC_AUTH_METHOD_TYPE_NONE"
  version                     = "OIDC_VERSION_1_0"
  clock_skew                  = "0s"
  dev_mode                    = false
  access_token_type           = "OIDC_TOKEN_TYPE_JWT"
  access_token_role_assertion = false
  id_token_role_assertion     = true
  id_token_userinfo_assertion = true
  additional_origins          = []
}

# Role Grant
resource "zitadel_project_role" "oidc_required_role" {
  org_id       = zitadel_org.website.id
  project_id   = zitadel_project.website.id
  role_key     = var.oidc_required_role
  display_name = var.oidc_required_role
  group        = var.oidc_required_role
}

resource "local_file" "zitadel_action" {
  filename        = "addGrant${var.oidc_required_role}.js"
  file_permission = "0664"

  content = templatefile("${path.module}/addGrant.js.tftpl", {
    oidc_required_role = var.oidc_required_role
    project_id         = zitadel_project.website.id
  })
}

resource "zitadel_action" "assign_role" {
  org_id          = zitadel_org.website.id
  name            = "addGrant${var.oidc_required_role}"
  script          = local_file.zitadel_action.content
  timeout         = "10s"
  allowed_to_fail = true
}

resource "zitadel_trigger_actions" "post_creation_trigger" {
  org_id       = zitadel_org.website.id
  flow_type    = "FLOW_TYPE_EXTERNAL_AUTHENTICATION"
  trigger_type = "TRIGGER_TYPE_POST_CREATION"
  action_ids = [
    zitadel_action.assign_role.id
  ]
}
