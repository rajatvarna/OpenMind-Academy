# This file defines the variables used in the Terraform configuration.
# Centralizing variables makes the infrastructure code more maintainable and reusable.

variable "gcp_project_id" {
  description = "The GCP Project ID to deploy resources into."
  type        = string
  default     = "free-education-platform-demo"
}

variable "gcp_region" {
  description = "The GCP region where resources will be created."
  type        = string
  default     = "us-central1"
}

variable "gke_cluster_name" {
  description = "The name for the Google Kubernetes Engine (GKE) cluster."
  type        = string
  default     = "edu-platform-gke-cluster"
}

variable "db_instance_name" {
  description = "The name for the Cloud SQL PostgreSQL instance."
  type        = string
  default     = "edu-platform-postgres-instance"
}

variable "db_name" {
  description = "The name of the main database."
  type        = string
  default     = "platform_db"
}

variable "db_user" {
  description = "The username for the main database user."
  type        = string
  default     = "platform_user"
  # In a real-world scenario, the password should be managed via a secret manager,
  # not a hardcoded variable.
}

variable "video_storage_bucket_name" {
  description = "The name for the GCS bucket to store user-generated videos. Must be globally unique."
  type        = string
  default     = "edu-platform-video-storage-bucket"
}
