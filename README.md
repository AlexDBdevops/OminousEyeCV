# OmnimousEyeCV

CV web en Go (generador estático) publicado en Cloudflare Pages. Infraestructura con Terraform (estado en HCP Terraform) y CI/CD con GitHub Actions. Coste: 0.

## Qué se edita
- `content/content.json`: textos ES/EN, experiencia, certificaciones y barras de nivel.
- `content/theme.txt`: tema activo (`cyan` o `red`); colores en `content/themes/*.json`.
- `static/`, `templates/`: diseño, comportamiento y cabeceras de seguridad (`static/headers.txt`).

Local: `go run .` genera `dist/` (o `THEME=red go run .`).

## Flujo
- Pull request: compila, `terraform plan` y vista previa en una URL `*.pages.dev` propia de la rama.
- Merge a `main`: `terraform apply` (proyecto de Pages en `infra/`) y publicación en producción.

Configuración en GitHub (Settings → Secrets and variables → Actions):
- Secretos: `CLOUDFLARE_API_TOKEN` (permiso Cloudflare Pages: Edit), `TF_API_TOKEN` (token de HCP Terraform).
- Variables: `CLOUDFLARE_ACCOUNT_ID`, `TF_CLOUD_ORGANIZATION`.
El workspace `omnimous-eye-cv` de HCP Terraform debe usar ejecución **Local**. Sin credenciales el workflow solo compila.
