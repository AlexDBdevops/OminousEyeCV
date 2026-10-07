// Genera static/og.png (1200x630) con el estilo de la v2, y los iconos apple-touch-icon / favicon-32.
// Uso: go run . && NODE_PATH=$(npm root -g) node tools/render-assets.js
const {chromium}=require('playwright');const fs=require('fs');const path=require('path');
const root=path.resolve(__dirname,'..'),dist=path.join(root,'dist');
(async()=>{
  const html=fs.readFileSync(path.join(dist,'index.html'),'utf8');
  const svg=html.match(/<svg id="art"[\s\S]*?<\/svg>/)[0];
  const banner=html.match(/<div class="banner bt" data-lang="es">([\s\S]*?)<\/div>/)[1];
  const og=`<!doctype html><html lang="es"><head><meta charset="utf-8">
<link rel="stylesheet" href="theme.css"><link rel="stylesheet" href="style.css"><link rel="stylesheet" href="variant-a.css">
<style>
html,body{margin:0;width:1200px;height:630px;overflow:hidden}
body.v-a{display:grid;grid-template-columns:540px 1fr;align-items:center;gap:28px;padding:0 56px 0 36px;box-sizing:border-box;
 background:radial-gradient(ellipse at 22% 50%,var(--bg2) 0,var(--bg) 72%)}
#art{width:540px!important;height:540px!important;overflow:visible}
#eyeball{animation:none}
.t{min-width:0;display:grid;gap:10px}
.p{font-family:var(--mono);font-size:19px;color:var(--bone)}.p span{color:var(--a)}
.t h1{font-family:var(--display);font-weight:800;font-size:58px;line-height:.95;letter-spacing:-.02em;margin:6px 0 0;color:var(--bone)}
.r{font-family:var(--mono);font-size:20px;color:var(--a3)}
.e{font-family:var(--mono);font-size:17px;color:var(--bone);margin-top:6px}.e span{color:var(--a)}
.b{color:var(--a3)}.b svg{display:block;width:100%;height:auto;fill:currentColor;filter:blur(.35px) drop-shadow(0 0 3px color-mix(in srgb,var(--a3) 70%,transparent)) drop-shadow(0 0 10px color-mix(in srgb,var(--a3) 40%,transparent))}
.k{font-family:var(--mono);font-size:18px;color:var(--bone);opacity:.85;margin-top:4px}
.u{font-family:var(--mono);font-size:18px;color:var(--mut)}
</style></head><body class="v-a">${svg}
<div class="t"><div class="p"><span>alejandro@ops:~$</span> ./AlexDBdevopsCV.sh</div>
<h1>Alejandro<br>Díaz Benjumea</h1><div class="r">&gt; sysadmin · devops</div>
<div class="e"><span>alejandro@ops:~$</span> echo "\\: I E N V E N | DOS"</div><div class="b">${banner}</div>
<div class="k">Linux · Docker · Kubernetes · GitOps</div><div class="u">alejandro-diaz-benjumea.pages.dev</div></div></body></html>`;
  const tmp=path.join(dist,'_og.html');fs.writeFileSync(tmp,og);
  const b=await chromium.launch({executablePath:'/opt/pw-browsers/chromium'});
  const p=await b.newPage({viewport:{width:1200,height:630}});
  await p.goto('file://'+tmp);await p.evaluate(()=>document.fonts.ready);
  await p.evaluate(()=>{const d=document.getElementById('disc');d.style.transition='none';d.style.transform='rotate(-90deg)';
    document.querySelectorAll('.sector')[0].classList.add('on');document.getElementById('iris').setAttribute('transform','translate(30 6)')});
  await p.waitForTimeout(300);
  await p.screenshot({path:path.join(root,'static/og.png')});
  fs.unlinkSync(tmp);
  const fav=fs.readFileSync(path.join(root,'static/favicon.svg'),'utf8');
  for(const [n,s] of [['apple-touch-icon.png',180],['favicon-32.png',32]]){
    const q=await b.newPage({viewport:{width:s,height:s}});
    await q.setContent(`<html><body style="margin:0;background:#07060d">${fav.replace('<svg ','<svg width="'+s+'" height="'+s+'" ')}</body></html>`);
    await q.screenshot({path:path.join(root,'static',n)});
  }
  await b.close();console.log('ok');
})();
