# This file defines the Google Cloud Storage (GCS) bucket
# for storing large media files, primarily the generated videos.

resource "google_storage_bucket" "video_storage" {
  # Append a random suffix to the bucket name to ensure global uniqueness.
  name          = "${var.video_storage_bucket_name}-${random_id.suffix.hex}"
  location      = var.gcp_region
  force_destroy = false # Set to true only for non-production environments

  # Uniform bucket-level access is recommended for simpler access control.
  uniform_bucket_level_access = true

  # Set up a lifecycle rule to transition older, less-accessed videos to
  # cheaper storage classes or delete them after a certain period.
  # This is a cost-optimization measure.
  lifecycle_rule {
    condition {
      age = 365 # Days
    }
    action {
      type = "SetStorageClass"
      storage_class = "NEARLINE"
    }
  }

  lifecycle_rule {
    condition {
      age = 730 # 2 years
    }
    action {
      type = "Delete"
    }
  }

  # Enable versioning to protect against accidental overwrites or deletions.
  versioning {
    enabled = true
  }
}

# Output the name of the created GCS bucket
output "video_storage_bucket_name_output" {
  value = google_storage_bucket.video_storage.name
}
