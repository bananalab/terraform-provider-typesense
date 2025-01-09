terraform {
  required_providers {
    typesense = {
      source = "bananalab/terraform/typesense"
    }
  }
}

data "typesense_cluster" "example" {
  id = "<cluster-id>"
}
