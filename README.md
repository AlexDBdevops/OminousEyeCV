# OmnimousEyeCV

CV web en Go (generador estático) publicado en Cloudflare Pages con GitHub Actions. Coste: 0.

## Qué se edita
- `content/content.json`: textos ES/EN, experiencia, certificaciones y barras de nivel.
- `content/theme.txt`: tema activo (`cyan` o `red`); colores en `content/themes/*.json`.
- `static/`, `templates/`: diseño, comportamiento y cabeceras de seguridad (`static/headers.txt`).

Local: `go run .` genera `dist/` (o `THEME=red go run .`).

## Flujo
- Pull request: compila y publica una vista previa en una URL `*.pages.dev` propia de la rama.
- Merge a `main`: publica en producción (`omnimous-eye-cv.pages.dev`).

Necesita en GitHub el secreto `CLOUDFLARE_API_TOKEN` (permiso Cloudflare Pages: Edit) y la variable `CLOUDFLARE_ACCOUNT_ID`.
Sin ellos el workflow solo compila.
