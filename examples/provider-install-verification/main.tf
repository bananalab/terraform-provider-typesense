terraform {
  required_providers {
    typesense = {
      source = "bananalab/terraform/typesense"
    }
  }
}

provider "typesense" {
}

data "typesense_cluster" "example" {
  id = "<cluster-id>"
}

output "example_cluster" {
  value = data.typesense_cluster.example
}
