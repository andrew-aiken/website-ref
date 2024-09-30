data "archive_file" "this" {
  type        = "zip"
  output_path = "function.zip"
  source_file = "files/index.js"
}
