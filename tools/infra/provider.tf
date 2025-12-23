terraform {
  required_version = "1.14.1"
  required_providers {
    keycloak = {
      source  = "keycloak/keycloak"
      version = "5.6.0"
    }
  }
}

provider "keycloak" {
  client_id = "admin-cli"
  url       = "http://auth:18080"
}
