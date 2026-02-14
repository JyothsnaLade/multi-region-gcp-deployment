output "primary_service_url" {
  description = "Primary region service URL"
  value       = google_cloud_run_v2_service.primary.uri
}

output "secondary_service_url" {
  description = "Secondary region service URL"
  value       = google_cloud_run_v2_service.secondary.uri
}

output "artifact_registry" {
  description = "Artifact Registry repository"
  value       = "${var.primary_region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.main.repository_id}"
}
