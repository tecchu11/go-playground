terraform {
  required_version = "1.11.0"
  required_providers {
    keycloak = {
      source  = "keycloak/keycloak"
      version = "5.1.1"
    }
  }
}

provider "keycloak" {
  client_id = "admin-cli"
  url       = "http://auth:18080"
}
