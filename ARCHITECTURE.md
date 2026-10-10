# OminousEyeCV · Cómo funciona

CV web de Alejandro Díaz Benjumea: **https://alejandro-diaz-benjumea.pages.dev**

Generador estático en Go, publicado en Cloudflare Pages, con la infraestructura en Terraform (estado en HCP Terraform) y CI/CD en GitHub Actions. Coste de alojamiento: 0.

## 1. Estructura del repositorio

```
content/content.json        ← todo el texto: experiencia, certificaciones, barras de nivel, "cómo está hecha"
content/theme.txt           ← tema activo (cyan o red)
content/themes/*.json       ← colores de cada tema
templates/index.html.tmpl   ← estructura HTML de la página
static/style.css            ← diseño
static/app.js               ← ojo, disco, botón ES/EN, barras animadas
static/headers.txt          ← cabeceras de seguridad (plantilla de _headers)
static/icons/*.svg          ← iconos de tecnologías (Simple Icons CC0, Devicon MIT y propios)
static/fonts/               ← fuentes autoalojadas (Syne, IBM Plex Sans/Mono, OFL)
static/variant-a.css        ← diseño de la v2 (estilo consola)
main.go                     ← generador: orquesta el build y escribe dist/
content.go                  ← tipos de content.json y su carga
disc.go                     ← geometría del disco (sectores, arcos, marcas)
render.go                   ← plantilla HTML, funciones e iconos en línea
headers.go                  ← hashes CSP de los scripts, _headers y theme.css
banner.go                   ← cartel BIENVENIDOS / WELCOME en SVG de píxeles
generator_test.go           ← tests del generador (go test ./...)
tools/build_single.sh       ← HTML autocontenido para vistas previas
tools/render-assets.js      ← og.png y favicons
tools/deploy/               ← Wrangler fijado a una versión exacta (package.json + package-lock.json)
tools/ci/                   ← zizmor (requirements.txt con hashes) y govulncheck (go.mod) fijados
docs/                       ← GIF y captura del README
VERSION                     ← versión publicada (v3, v3.1, v3.1.1…)
LICENSE                     ← MIT (fuentes OFL, iconos CC0/MIT)
infra/main.tf               ← Terraform: proyecto de Cloudflare Pages
.github/workflows/deploy.yml← CI/CD
.github/workflows/tag.yml   ← crea la etiqueta de versión cuando cambia VERSION en main
.github/dependabot.yml      ← actualizaciones automáticas de dependencias
```

El contenido está separado del diseño: para cambiar un texto o un porcentaje solo se toca `content.json`; para cambiar el color, una palabra en `theme.txt`.

## 2. El generador (`main.go` y compañía)

`go run .` hace lo siguiente:

1. Lee `content.json`, con cada texto en español e inglés.
2. Calcula la geometría del disco: un sector por sección (ahora dos de 180°), los arcos para el texto curvado, las marcas del borde y los radios del iris.
3. Rellena la plantilla HTML. Cada texto se escribe dos veces (`data-lang="es"` y `data-lang="en"`) y los iconos SVG se insertan en línea.
4. Escribe en `dist/`: `index.html`, `style.css`, `app.js`, `theme.css` (colores del tema elegido) y `_headers` (seguridad).
5. Calcula el hash SHA-256 del script en línea del HTML y lo añade a la Content-Security-Policy, de modo que el navegador solo ejecuta ese script y ninguno inyectado.

El resultado es una web estática: archivos fijos, sin servidor ni base de datos.

## 3. En el navegador (`app.js`)

- **Idioma:** el HTML trae los dos idiomas y el CSS oculta uno según `<html lang>`. Por defecto, inglés; el botón cambia y el navegador lo recuerda.
- **Ojo:** el movimiento del ratón o un toque fija una posición objetivo para el iris; en cada fotograma el iris avanza un 16 % hacia ella (movimiento suave). El párpado recorta el iris. Parpadeo cada 6 s con CSS. Con "reducir movimiento" activado, el iris sigue al puntero sin interpolación.
- **Disco:** pulsar un sector, un botón o arrastrar el anillo lo gira por el camino más corto hasta dejar arriba la sección elegida y muestra su bloque de experiencia.
- **Barras:** se rellenan con animación al entrar en pantalla.
- **Arranque (v2):** la primera visita de cada sesión empieza con solo el prompt; se escribe `./AlexDBdevopsCV.sh`, aparece un botón `↵ Enter` (o arranca solo a los 3 s) y el CV se va mostrando por partes hasta `finish in 3.3s`. Cualquier tecla o clic lo salta; si el JavaScript fallara, todo aparece a los 7 s.
- **Cartel:** `echo` con BIENVENIDOS / WELCOME dibujado por Go como SVG de píxeles, con líneas de barrido y brillo.
- **Pie:** commit y fecha reales del build.

## 4. Infraestructura (`infra/main.tf`)

Terraform gestiona el proyecto de Cloudflare Pages `alejandro-diaz-benjumea` (que da el dominio `.pages.dev`). El estado se guarda en HCP Terraform, workspace `ominous-eye-cv`, en modo de ejecución **Local**: Terraform se ejecuta en GitHub Actions y HCP solo guarda y bloquea el estado. Cambiar el nombre del proyecto en Terraform sustituye el proyecto (borra el antiguo y crea el nuevo).

**Terraform vs HCP Terraform:** Terraform es la herramienta que compara el código con lo que existe en Cloudflare y aplica los cambios. HCP Terraform es el servicio que guarda el archivo de estado de forma remota, lo bloquea durante cada ejecución y conserva su historial.

## 5. CI/CD (`.github/workflows/deploy.yml`)

Se lanza con cualquier push (a `main` o a otra rama) o a mano (Actions → deploy → Run workflow):

1. **Checkout y Go.**
2. **Calidad:** `gofmt` (formato), `go vet` y `go test`. Si algo falla, se para y la web no cambia.
   **Seguridad:** `zizmor` revisa los workflows (inyecciones, permisos, acciones sin fijar) y `govulncheck` busca vulnerabilidades conocidas que afecten al generador. Las dos herramientas tienen la versión y los hashes fijados en `tools/ci`.
3. **Compilar:** `go run .` genera `dist/`.
4. **Comprobar credenciales:** si falta algún secreto, solo compila.
5. **Terraform init + fmt + validate.**
6. **Otras ramas:** `terraform plan` (muestra qué cambiaría, sin tocar nada). **main:** `terraform apply`.
7. **Wrangler** (CLI de Cloudflare) sube `dist/` al proyecto de Pages: a producción en `main` y, en cualquier otra rama, a una URL de vista previa propia (`<rama>.alejandro-diaz-benjumea.pages.dev`). No hace falta abrir un pull request.
8. **Redirección `www`:** en `main` se despliega también una rama `www` que solo contiene un `_redirects` con un 301 a la dirección principal, porque en `pages.dev` el prefijo `www.` se interpreta como nombre de rama.
9. **Limpieza (`cleanup.yml`, llamado como último job de `deploy.yml`):** tras cada despliegue, al borrar una rama o a mano, `tools/cleanup-deployments.sh` borra con la API de Cloudflare los despliegues que sobran: deja los 3 últimos de producción (para poder volver atrás), el último de `www` y el de cada rama que siga existiendo. Al borrar la rama de una prueba, su vista previa desaparece sola.

Versionado: los hitos se marcan cambiando `VERSION` en `main`; el workflow `tag` crea la etiqueta. Los cambios pequeños (una skill, un nivel, una frase) suben el último dígito (v3.2 → v3.2.1); los más grandes, el del medio (v3.2 → v3.3). Solo llevan versión los cambios en la web. Un push a `main` despliega si toca algo de la web o de su infraestructura; si solo cambian documentación, changelog, `VERSION`, `tools/` o los workflows `tag` y `cleanup`, se ejecutan los tests pero no se despliega (paso «Detect site changes» de `deploy.yml`).

Configuración en GitHub (Settings → Secrets and variables → Actions):

| Tipo | Nombre | Para qué |
|---|---|---|
| Secreto | `CLOUDFLARE_API_TOKEN` | Token con permiso solo de Cloudflare Pages: Edit |
| Secreto | `TF_API_TOKEN` | Token de HCP Terraform para el estado |
| Variable | `CLOUDFLARE_ACCOUNT_ID` | ID de la cuenta de Cloudflare |
| Variable | `TF_CLOUD_ORGANIZATION` | Organización de HCP Terraform |

## 6. Seguridad

- Cabeceras: CSP sin `unsafe-inline` para scripts (hash calculado en el build), HSTS, `X-Frame-Options: DENY`, COOP/CORP, `Permissions-Policy` (sin cámara, micrófono, ubicación ni pagos), `nosniff`, `no-referrer`.
- Workflow: acciones fijadas por SHA, `GITHUB_TOKEN` de solo lectura, checkout sin credenciales persistidas y tiempo límite.
- Credenciales fuera del código, con permisos mínimos.
- Denegación de servicio: Cloudflare filtra los ataques DDoS en todos los planes y el tráfico estático no se cobra, así que un ataque no genera coste.

## 7. Dependabot

Servicio integrado de GitHub configurado en `.github/dependabot.yml`. Cada mes revisa si hay versiones nuevas de las acciones del workflow, del proveedor de Terraform, de Wrangler (`tools/deploy`) y de zizmor y govulncheck (`tools/ci`), con un `cooldown` de 7 días: no propone una versión hasta que lleva una semana publicada, para no coger una recién comprometida antes de que se detecte. Si las hay, si las hay, abre un pull request con la actualización. Su rama ejecuta el workflow (formato, tests y compilación; GitHub no da los secretos a Dependabot, así que sin despliegue), así que se puede comprobar que todo sigue funcionando antes de fusionarlo. No fusiona nada por sí solo.

