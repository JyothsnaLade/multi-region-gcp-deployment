# Artifact Registry for container images
resource "google_artifact_registry_repository" "main" {
  location      = var.primary_region
  repository_id = "${var.service_name}-repo"
  description   = "Docker repository for ${var.service_name}"
  format        = "DOCKER"
}

# Cloud Run service in primary region
resource "google_cloud_run_v2_service" "primary" {
  name     = "${var.service_name}-${var.environment}-primary"
  location = var.primary_region

  template {
    containers {
      image = var.image

      env {
        name  = "REGION"
        value = var.primary_region
      }
      env {
        name  = "VERSION"
        value = "1.0.0"
      }
      env {
        name  = "ENVIRONMENT"
        value = var.environment
      }
    }

    scaling {
      min_instance_count = var.min_instances
      max_instance_count = var.max_instances
    }
  }
}

# Make service publicly accessible
resource "google_cloud_run_v2_service_iam_member" "primary_public" {
  location = google_cloud_run_v2_service.primary.location
  name     = google_cloud_run_v2_service.primary.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Cloud Run service in secondary region
resource "google_cloud_run_v2_service" "secondary" {
  name     = "${var.service_name}-${var.environment}-secondary"
  location = var.secondary_region

  template {
    containers {
      image = var.image

      env {
        name  = "REGION"
        value = var.secondary_region
      }
      env {
        name  = "VERSION"
        value = "1.0.0"
      }
      env {
        name  = "ENVIRONMENT"
        value = var.environment
      }
    }

    scaling {
      min_instance_count = var.min_instances
      max_instance_count = var.max_instances
    }
  }
}

# Make service publicly accessible
resource "google_cloud_run_v2_service_iam_member" "secondary_public" {
  location = google_cloud_run_v2_service.secondary.location
  name     = google_cloud_run_v2_service.secondary.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
