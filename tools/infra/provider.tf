terraform {
  required_version = "1.12.2"
  required_providers {
    keycloak = {
      source  = "keycloak/keycloak"
      version = "5.4.0"
    }
  }
}

provider "keycloak" {
  client_id = "admin-cli"
  url       = "http://auth:18080"
}
