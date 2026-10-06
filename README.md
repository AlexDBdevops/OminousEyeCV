# OmnimousEyeCV

CV web en Go (generador estático) desplegado en S3 + CloudFront con Terraform y GitHub Actions.

## Qué se edita
- `content/content.json`: textos ES/EN, experiencia, competencias.
- `content/theme.txt`: tema activo (`red` o `cyan`). Los colores están en `content/themes/*.json`.
- `static/`, `templates/`: diseño y comportamiento.

Local: `go run .` genera `dist/` (o `THEME=cyan go run .`).

## Flujo
1. Rama + pull request: se compila el sitio y se ejecuta `terraform plan` (rol de solo lectura).
2. Merge a `main`: `terraform apply` + subida a S3 + invalidación de CloudFront (entorno `production`).

La puesta en marcha inicial está en `infra/bootstrap.yaml`. No hay claves de AWS en el repositorio.
