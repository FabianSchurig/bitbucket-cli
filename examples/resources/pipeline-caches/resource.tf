resource "bitbucket_pipeline_caches" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  cache_uuid = "{cache-uuid}"
}
