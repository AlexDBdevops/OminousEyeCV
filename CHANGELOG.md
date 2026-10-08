# Changelog

## v3.0.1 · bugfix

- **Reload left the page scrolled down during the boot.** After scrolling and reloading without cache in a new session, the browser restored the old scroll position while the simulated `./AlexDBdevopsCV.sh` ran at the top, so the visitor saw an empty page that looked like it was still loading. The page now starts at the prompt, stays there while the boot runs, and scrolling with the wheel or by touch skips the boot. Before / after: [docs/bugfix-reload-scroll.mp4](docs/bugfix-reload-scroll.mp4).
- Skills: same spacing between every skill when the columns stack on narrow screens.
- Experience: the DevSecOps block lines up with the role header, with dividers between columns.
- `www.` redirects to the main address; old Cloudflare Pages deployments are cleaned up automatically.

## v3

Clean history. Bilingual (EN/ES) CV site generated with Go and deployed to Cloudflare Pages with Terraform (state in HCP Terraform) and GitHub Actions.
