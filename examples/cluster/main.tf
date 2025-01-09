terraform {
  required_providers {
    typesense = {
      source = "bananalab/terraform/typesense"
    }
  }
}

resource "typesense_cluster" "example" {
  memory            = "0.5_gb"
  vcpu              = "2_vcpus_1_hr_burst_per_day"
  region            = "oregon"
  name              = "example"
  high_availability = "no"
}

output "typesense_cluster" {
  value = typesense_cluster.example
}

resource "typesense_cluster_api_keys" "example" {
  cluster_id = typesense_cluster.example.id
}

output "typesense_cluster_api_keys" {
  value       = typesense_cluster_api_keys.example
  description = "Admin key"
  sensitive   = true
}
