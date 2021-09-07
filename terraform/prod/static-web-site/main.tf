module "static-web-site" {
  source       = "https://artifactory.kindredgroup.com/artifactory/terraform-modules-local/pes/terraform-static-website/terraform-static-website-0.0.5.tar.gz"
  bucket       = "${var.project_name}"
  content_path = path.module
  policy       = {
    "Effect" = "Allow"
    "Principal" = "*"
    "Action" = "s3:GetObject"
  }
  policy_allowed_source_ip =[
    "0.0.0.0"
    ]
}