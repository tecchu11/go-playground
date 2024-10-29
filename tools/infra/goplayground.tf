locals {
  app_name = "go_playground"
}

resource "keycloak_realm" "go_playground" {
  realm                    = local.app_name
  enabled                  = true
  display_name             = local.app_name
  registration_allowed     = true
  reset_password_allowed   = false
  remember_me              = true
  login_with_email_allowed = false
  verify_email             = false
  access_token_lifespan    = "30m"
}

resource "keycloak_openid_client" "go_playground_frontend" {
  realm_id              = keycloak_realm.go_playground.id
  client_id             = "frontend"
  description           = "client for frontend"
  standard_flow_enabled = true
  valid_redirect_uris = [
    "http://localhost:3000/*"
  ]
  web_origins = [
    "http://localhost:3000"
  ]
  access_type = "PUBLIC"
  login_theme = "keycloak.v2"
}
