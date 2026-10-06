#!/bin/sh
# uso: ./build_single.sh <tema> <salida.html>
set -e
THEME=$1 go run . >/dev/null
python3 - "$2" <<'PY'
import re,sys
d='dist/'
h=open(d+'index.html').read()
h=h.replace('<link rel="stylesheet" href="theme.css">','<style>'+open(d+'theme.css').read()+'</style>')
h=h.replace('<link rel="stylesheet" href="style.css">','<style>'+open(d+'style.css').read()+'</style>')
h=h.replace('<script src="app.js"></script>','<script>'+open(d+'app.js').read()+'</script>')
h=re.sub(r'<title>.*?</title>','<title>Omnimous Eye CV</title>',h)
open(sys.argv[1],'w').write(h)
PY
