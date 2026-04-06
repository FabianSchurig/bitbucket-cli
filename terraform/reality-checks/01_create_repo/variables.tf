variable "workspace" {
  description = "Bitbucket workspace slug (set via TF_VAR_workspace or in .env at repository root)"
  type = string
}

variable "hook_url" {
  description = "Webhook delivery URL (default points to example.com)"
  type = string
  default = "https://example.com/hook"
}

variable "create_branch_restriction" {
  description = "Whether to create a branch restriction during the reality-check. Defaults to false because the provider/schema may not support all kinds."
  type        = bool
  default     = true
}

variable "create_hook" {
  description = "Whether to create a webhook during the reality-check. Defaults to false because the provider currently doesn't expose webhook body fields."
  type        = bool
  default     = true
}
