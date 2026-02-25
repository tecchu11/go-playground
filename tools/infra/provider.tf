terraform {
  required_version = "1.14.6"
  required_providers {
    keycloak = {
      source  = "keycloak/keycloak"
      version = "5.7.0"
    }
  }
}

provider "keycloak" {
  client_id = "admin-cli"
  url       = "http://auth:18080"
}
