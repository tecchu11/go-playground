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
  account_theme            = "keycloak.v3"
  admin_theme              = "keycloak.v2"
  email_theme              = "keycloak"
  login_theme              = "keycloak.v2"
}

resource "keycloak_openid_client" "frontend" {
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

resource "keycloak_openid_client" "backend" {
  realm_id    = keycloak_realm.go_playground.id
  client_id   = "backend"
  description = "client for backend"
  access_type = "BEARER-ONLY"
}

resource "keycloak_openid_client_scope" "backend_access" {
  realm_id               = keycloak_realm.go_playground.id
  name                   = "backend-access"
  description            = "When requested, this scope will map a user's group memberships to a claim"
  include_in_token_scope = true
}

resource "keycloak_openid_audience_protocol_mapper" "frontend" {
  realm_id                 = keycloak_realm.go_playground.id
  client_scope_id          = keycloak_openid_client_scope.backend_access.id
  included_client_audience = "backend"
  name                     = "front-backend-mapper"
  depends_on               = [keycloak_openid_client_scope.backend_access]
}

resource "keycloak_openid_client_optional_scopes" "frontend" {
  realm_id  = keycloak_realm.go_playground.id
  client_id = keycloak_openid_client.frontend.id
  optional_scopes = [
    keycloak_openid_client_scope.backend_access.name,
  ]
}
