variable "workspace" {
  description = "Bitbucket workspace slug (set via TF_VAR_workspace or in .env at repository root)"
  type        = string
}

variable "feature_prefix" {
  description = <<-EOT
    Prefix for the "feature" branch type. Change it on a second apply
    (e.g. -var feature_prefix=feature2/) to prove the branching-model update is
    detected — this exercises the requestBody/typed-fields fix (the update
    operations had HasBody=false until the schema was re-enriched).
  EOT
  type        = string
  default     = "feature/"
}

variable "manage_project_model" {
  description = "Also create a project and configure its branching model."
  type        = bool
  default     = true
}
