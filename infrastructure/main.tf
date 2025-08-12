# This is the main configuration file for the project's infrastructure.
# It sets up the Google Cloud provider and defines the backend for storing
# the Terraform state.

# Configure the Google Cloud provider
provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

# Configure the Terraform backend to store the state file in a GCS bucket.
# This is crucial for collaboration and state management in a team environment.
# Note: This bucket must be created manually before running 'terraform init'.
terraform {
  backend "gcs" {
    bucket = "education-platform-tf-state-bucket" # A unique name for the GCS bucket
    prefix = "terraform/state"
  }
}

# Define a random suffix to ensure resource names are unique
resource "random_id" "suffix" {
  byte_length = 4
}
