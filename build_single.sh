#!/bin/sh
# Genera una página única autocontenida (CSS, JS y fuentes en línea) para revisión.
# uso: ./build_single.sh <tema> <salida.html> [variante]
set -e
THEME=$1 VARIANT=${3:-$(cat content/variant.txt)} go run . >/dev/null
python3 - "$2" "${3:-$(cat content/variant.txt)}" <<'PY'
import re,sys,base64,os
d='dist/'; out,var=sys.argv[1],sys.argv[2]
h=open(d+'index.html').read()
def css(name):
    s=open(d+name).read()
    return re.sub(r'url\((fonts/[^)]+)\)',lambda m:'url(data:font/woff;base64,'+base64.b64encode(open(d+m.group(1),'rb').read()).decode()+')',s)
h=h.replace('<link rel="stylesheet" href="theme.css">','<style>'+css('theme.css')+'</style>')
h=h.replace('<link rel="stylesheet" href="style.css">','<style>'+css('style.css')+'</style>')
h=h.replace(f'<link rel="stylesheet" href="variant-{var}.css">','<style>'+css(f'variant-{var}.css')+'</style>')
h=h.replace('<script src="app.js"></script>','<script>'+open(d+'app.js').read()+'</script>')
title='Omnimous Eye CV' if var=='a' and 'v2' not in out else None
h=re.sub(r'<title>.*?</title>',f'<title>{os.environ.get("TITLE","Omnimous Eye CV")}</title>',h)
open(out,'w').write(h)
PY
