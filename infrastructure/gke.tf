# This file defines the Google Kubernetes Engine (GKE) cluster.
# GKE will orchestrate and manage our containerized microservices.

resource "google_container_cluster" "primary" {
  name     = var.gke_cluster_name
  location = var.gcp_region

  # We start with a single node pool. More specialized node pools can be added later.
  # e.g., a high-CPU pool for video processing or a high-memory pool for AI services.
  initial_node_count = 1
  remove_default_node_pool = true

  # Using a cost-effective, general-purpose machine type to start.
  # This can be scaled up as the platform grows.
  node_config {
    machine_type = "e2-medium"
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

resource "google_container_node_pool" "primary_nodes" {
  name       = "primary-node-pool"
  cluster    = google_container_cluster.primary.name
  location   = google_container_cluster.primary.location
  node_count = 2

  node_config {
    machine_type = "e2-medium"
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

# Output the GKE cluster name and endpoint for connecting with kubectl
output "gke_cluster_name" {
  value = google_container_cluster.primary.name
}

output "gke_cluster_endpoint" {
  value = google_container_cluster.primary.endpoint
}
