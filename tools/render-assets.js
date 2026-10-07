// Genera static/og.png (1200x630), static/apple-touch-icon.png y static/favicon-32.png.
// Uso: go run . && NODE_PATH=$(npm root -g) node tools/render-assets.js
const {chromium}=require('playwright');const fs=require('fs');const path=require('path');
const root=path.resolve(__dirname,'..');
(async()=>{
  const html=fs.readFileSync(path.join(root,'dist/index.html'),'utf8');
  const svg=html.match(/<svg id="art"[\s\S]*?<\/svg>/)[0].replace('id="art"','id="art" width="560" height="560"');
  const css=fs.readFileSync(path.join(root,'dist/theme.css'),'utf8')+fs.readFileSync(path.join(root,'dist/style.css'),'utf8');
  const og=`<!doctype html><html lang="en"><head><meta charset="utf-8"><style>${css}
  html,body{margin:0;width:1200px;height:630px;overflow:hidden}
  body{background:radial-gradient(ellipse at 25% 50%,var(--bg2) 0,var(--bg) 70%);display:flex;align-items:center;gap:36px;padding:0 64px 0 40px;box-sizing:border-box;font-family:system-ui,-apple-system,"Segoe UI",Roboto,sans-serif}
  #eyeball{animation:none}#art{width:520px!important;height:520px!important;flex:none;overflow:visible}.seglabel{font-size:15px}
  .t{min-width:0}.t h1{font-size:60px;line-height:1.02;margin:0;color:var(--bone)}
  .t p{margin:16px 0 0;color:var(--a);letter-spacing:.22em;text-transform:uppercase;font-size:22px}
  .t .u{color:var(--mut);letter-spacing:.02em;text-transform:none;font-size:21px;margin-top:28px;white-space:nowrap}
  .t .k{color:var(--bone);letter-spacing:.02em;text-transform:none;font-size:22px;margin-top:22px;opacity:.85}
  </style></head><body>${svg}<div class="t"><h1>Alejandro<br>Díaz Benjumea</h1><p>Sysadmin · DevOps</p>
  <p class="k">Kubernetes · CI/CD · GitOps · Terraform</p><p class="u">alejandro-diaz-benjumea.pages.dev</p></div></body></html>`;
  const b=await chromium.launch({executablePath:'/opt/pw-browsers/chromium'});
  const p=await b.newPage({viewport:{width:1200,height:630}});
  await p.setContent(og);
  await p.evaluate(()=>{document.getElementById('disc').style.transition='none';document.getElementById('disc').style.transform='rotate(-90deg)';
    document.querySelectorAll('.sector')[0].classList.add('on');document.getElementById('iris').setAttribute('transform','translate(26 4)')});
  await p.screenshot({path:path.join(root,'static/og.png')});
  const fav=fs.readFileSync(path.join(root,'static/favicon.svg'),'utf8');
  for(const [n,s] of [['apple-touch-icon.png',180],['favicon-32.png',32]]){
    const q=await b.newPage({viewport:{width:s,height:s}});
    await q.setContent(`<html><body style="margin:0;background:#07060d">${fav.replace('<svg ','<svg width="'+s+'" height="'+s+'" ')}</body></html>`);
    await q.screenshot({path:path.join(root,'static',n),omitBackground:false});
  }
  await b.close();console.log('ok');
})();
