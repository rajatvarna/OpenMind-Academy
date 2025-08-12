# This file defines the Cloud SQL for PostgreSQL instance.
# A managed database is used to reduce operational overhead.

resource "google_sql_database_instance" "postgres" {
  name             = var.db_instance_name
  database_version = "POSTGRES_14"
  region           = var.gcp_region

  settings {
    # Using a small machine type to start, suitable for development/staging.
    # This should be scaled up for production.
    tier = "db-g1-small"
  }
}

# Define the main database within the PostgreSQL instance
resource "google_sql_database" "main_db" {
  name     = var.db_name
  instance = google_sql_database_instance.postgres.name
}

# Define the primary user for the database
resource "google_sql_user" "main_user" {
  name     = var.db_user
  instance = google_sql_database_instance.postgres.name
  # The password should be generated and stored in Google Secret Manager.
  # For this example, we use a random password.
  password = random_password.db_password.result
}

# Generate a random password for the database user.
# In a real setup, this would be managed by a secrets manager.
resource "random_password" "db_password" {
  length  = 16
  special = true
}

# Output the database connection details.
# In a real environment, services running on GKE would use the Cloud SQL Proxy
# for secure connections, rather than public IP.
output "db_instance_connection_name" {
  value = google_sql_database_instance.postgres.connection_name
}

output "db_name_output" {
  value = google_sql_database.main_db.name
}

output "db_user_output" {
  value = google_sql_user.main_user.name
}

output "db_password_output" {
  description = "The generated password for the database user. Store this securely."
  value       = random_password.db_password.result
  sensitive   = true
}
