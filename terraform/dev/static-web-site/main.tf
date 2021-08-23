module "static-web-site" {
  source       = "https://artifactory.kindredgroup.com/artifactory/terraform-modules-local/pes/terraform-static-website/terraform-static-website-0.0.5.tar.gz"
  bucket       = "dev.${var.project_name}"
  content_path = path.module
  policy       = {
    "Effect" = "Allow"
    "Principal" = "*"
    "Action" = "s3:GetObject"
  }
  policy_allowed_source_ip =[
    "85.184.100.0/24",
    "85.184.102.0/24",
    "185.9.100.0/24",
    "192.121.200.0/24",
    "185.83.149.0/24"
    ]
}
