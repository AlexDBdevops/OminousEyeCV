# Infraestructura del CV en Cloudflare, gestionada con Terraform.
# Estado remoto en HCP Terraform (plan gratuito). Organización y workspace llegan por
# variables de entorno: TF_CLOUD_ORGANIZATION y TF_WORKSPACE (ver el workflow).
terraform {
  required_version = ">= 1.6"
  cloud {}
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.0"
    }
  }
}

# El token se lee de la variable de entorno CLOUDFLARE_API_TOKEN.
provider "cloudflare" {}

variable "cloudflare_account_id" {
  type        = string
  description = "ID de la cuenta de Cloudflare"
}

variable "project_name" {
  type    = string
  default = "omnimous-eye-cv"
}

# Proyecto de Pages en modo "direct upload": GitHub Actions compila con Go y sube dist/.
resource "cloudflare_pages_project" "cv" {
  account_id        = var.cloudflare_account_id
  name              = var.project_name
  production_branch = "main"
}

output "url" {
  value = "https://${cloudflare_pages_project.cv.subdomain}"
}
