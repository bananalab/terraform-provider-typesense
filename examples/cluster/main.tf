terraform {
  required_providers {
    typesense = {
      source = "omarkhd.net/terraform/typesense"
    }
  }
}

provider "typesense" {
  key = "sUQKJv6AafWMSWseUyKFtaiY::7d32d883-942e-4baa-ab13-1a22baa0d97d"
}

resource "typesense_cluster" "example" {
  memory            = "0.5_gb"
  vcpu              = "2_vcpus_1_hr_burst_per_day"
  region            = "oregon"
  name              = "example"
  high_availability = "no"
}

resource "typesense_cluster_api_keys" "example" {
  cluster_id = typesense_cluster.example.id
}

output "typesense-admin-cluster-api-key" {
  value       = typesense_cluster_api_keys.example.admin_key
  description = "Admin key"
}

output "typesense-search-only-cluster-api-key" {
  value       = typesense_cluster_api_keys.example.search_only_key
  description = "Search Only key"
}
