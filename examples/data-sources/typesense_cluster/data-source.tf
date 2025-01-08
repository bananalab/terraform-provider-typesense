terraform {
  required_providers {
    typesense = {
      source = "omarkhd.net/terraform/typesense"
    }
  }
}
data "typesense_cluster" "example" {
  id = "s89j4uytxnbhomfcp"
}
