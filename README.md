# OminousEyeCV

[![deploy](https://github.com/AlexDBdevops/OminousEyeCV/actions/workflows/deploy.yml/badge.svg)](https://github.com/AlexDBdevops/OminousEyeCV/actions/workflows/deploy.yml)
[![license: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

My CV as a website: **https://alejandro-diaz-benjumea.pages.dev**

A bilingual (EN/ES) static site generated with Go, deployed to Cloudflare Pages with Terraform and GitHub Actions. Hosting cost: €0.

![Boot sequence](docs/boot.gif)

## How it works

```mermaid
flowchart LR
    dev[push / pull request] --> gha[GitHub Actions]
    gha --> q[gofmt · go vet · go test]
    q --> build[go run . → dist/]
    build --> tf[Terraform<br/>state in HCP Terraform]
    tf --> wr[Wrangler]
    wr --> cf[(Cloudflare Pages)]
    cf -->|main| prod[production]
    cf -->|PR| prev[preview URL per branch]
```

- **Pull request:** tests, build, `terraform plan` and a preview deploy on its own `*.pages.dev` URL.
- **Merge to `main`:** `terraform apply` and production deploy.
- **Releases:** bump `VERSION` on `main` and the `tag` workflow creates the git tag (v1, v2, v2.1…).

## Stack and why

| Piece | Why |
|---|---|
| **Go** static generator (`html/template`, `embed`) | One binary, no runtime dependencies, output is plain files. |
| **Cloudflare Pages** | Free static hosting with global CDN and DDoS protection included; static traffic is never billed. Chosen over AWS (S3 + CloudFront) for guaranteed zero cost. |
| **Terraform** + **HCP Terraform** | The Pages project is infrastructure as code; remote state with locking and history. |
| **GitHub Actions** | CI/CD with actions pinned by SHA, read-only token and no persisted credentials. |
| **Dependabot** | Monthly PRs for action and provider updates, validated by the same pipeline. |

Security: strict Content-Security-Policy (inline scripts allowed only by SHA-256 hash computed at build time), HSTS, `X-Frame-Options: DENY`, COOP/CORP, Permissions-Policy, self-hosted fonts and inline icons, so the page makes no third-party requests. Credentials live outside the repo with minimum permissions.

## Run locally

```sh
go run .            # builds dist/  (THEME=red go run . for the alternative theme)
go test ./...       # generator tests
TITLE="CV" tools/build_single.sh cyan "$PWD/single.html" a   # single self-contained HTML file
go run . && NODE_PATH=$(npm root -g) node tools/render-assets.js   # regenerate og.png and favicons (Playwright)
```

## Layout

```
content/content.json     all text (EN/ES), experience, skills and levels
content/theme.txt        active theme · content/themes/*.json colours
content/variant.txt      design variant (a = console style)
templates/               HTML template
static/                  CSS, JS, icons, fonts, security headers template
main.go                  build orchestration → dist/
content.go  disc.go  render.go  headers.go  banner.go   generator pieces
generator_test.go        tests
infra/                   Terraform (Cloudflare Pages project)
.github/workflows/       deploy (CI/CD) and tag (releases)
tools/                   single-file build and asset rendering
```

More detail (in Spanish) in [ARCHITECTURE.md](ARCHITECTURE.md).

## Credits

- Built with the assistance of [Claude Code](https://claude.com/claude-code), which wrote much of the code under my direction and review.
- Icons: [Simple Icons](https://simpleicons.org) (CC0) and [Devicon](https://devicon.dev) (MIT).
- Fonts: Syne and IBM Plex (SIL OFL 1.1).

Licensed under [MIT](LICENSE).
