(function(){
var root=document.documentElement,svg=document.getElementById('art'),iris=document.getElementById('iris'),
disc=document.getElementById('disc'),sectors=[].slice.call(document.querySelectorAll('.sector')),
pills=[].slice.call(document.querySelectorAll('.pill')),details=[].slice.call(document.querySelectorAll('.detail'));
var reduce=matchMedia('(prefers-reduced-motion: reduce)').matches;

// Eye tracking: el iris persigue al puntero con interpolación suave (rAF)
var tx=0,ty=0,cx0=0,cy0=0,raf=0;
function look(x,y){
  var r=svg.getBoundingClientRect(),cx=r.left+r.width/2,cy=r.top+r.height/2;
  var dx=x-cx,dy=y-cy,dist=Math.hypot(dx,dy),a=Math.atan2(dy,dx);
  var d=1-Math.exp(-dist/(r.width*0.22));          // satura rápido: responde aunque el ratón esté lejos
  tx=Math.cos(a)*62*d; ty=Math.sin(a)*34*d;        // en unidades del viewBox (600)
  if(!raf)raf=requestAnimationFrame(step);
}
function step(){
  var k=reduce?1:0.16;               // con movimiento reducido: sin interpolación, pero sigue al puntero
  cx0+=(tx-cx0)*k; cy0+=(ty-cy0)*k;
  iris.setAttribute('transform','translate('+cx0.toFixed(2)+' '+cy0.toFixed(2)+')');
  raf=(Math.abs(tx-cx0)>0.05||Math.abs(ty-cy0)>0.05)?requestAnimationFrame(step):0;
}
{
  addEventListener('pointermove',function(e){look(e.clientX,e.clientY)},{passive:true});
  addEventListener('mousemove',function(e){look(e.clientX,e.clientY)},{passive:true});
  addEventListener('pointerdown',function(e){look(e.clientX,e.clientY)},{passive:true});
  document.addEventListener('mouseleave',function(){tx=0;ty=0;if(!raf)raf=requestAnimationFrame(step)});
}

// Disc selection
var rot=0,cur=-1,N=sectors.length,SP=360/N;
function select(i){
  if(i===cur)return; cur=i;
  var target=-(i*SP+SP/2),delta=((target-rot)%360+540)%360-180; // shortest turn
  rot+=delta; disc.style.transform='rotate('+rot+'deg)';
  sectors.forEach(function(s,j){s.classList.toggle('on',j===i)});
  pills.forEach(function(p,j){p.classList.toggle('on',j===i)});
  details.forEach(function(d,j){d.classList.toggle('on',j===i)});
}
sectors.forEach(function(s){
  var i=+s.dataset.i;
  s.addEventListener('click',function(){select(i)});
  s.addEventListener('keydown',function(e){if(e.key==='Enter'||e.key===' '){e.preventDefault();select(i)}});
});
pills.forEach(function(p){p.addEventListener('click',function(){select(+p.dataset.i)})});

// Drag to spin: snap to nearest sector
var drag=null;
svg.addEventListener('pointerdown',function(e){
  var r=svg.getBoundingClientRect(),cx=r.left+r.width/2,cy=r.top+r.height/2,
  d=Math.hypot(e.clientX-cx,e.clientY-cy)/(r.width/2);
  if(d<0.68||d>1)return;
  drag={cx:cx,cy:cy,a:Math.atan2(e.clientY-cy,e.clientX-cx),r:rot,moved:false};
});
addEventListener('pointermove',function(e){
  if(!drag)return;
  var a=Math.atan2(e.clientY-drag.cy,e.clientX-drag.cx),da=(a-drag.a)*180/Math.PI;
  if(Math.abs(da)>4)drag.moved=true;
  if(drag.moved){disc.style.transition='none';disc.style.transform='rotate('+(drag.r+da)+'deg)';drag.cur=drag.r+da}
});
addEventListener('pointerup',function(){
  if(!drag)return;
  if(drag.moved){
    disc.style.transition='';
    var ang=(((-drag.cur)%360)+360)%360,i=Math.floor(ang/SP)%N;
    rot=drag.cur;cur=-1;select(i);
  }
  drag=null;
});
select(0);

// Language
document.getElementById('lang').addEventListener('click',function(){
  var l=root.lang==='es'?'en':'es';root.lang=l;
  try{localStorage.setItem('lang',l)}catch(e){}
});

// Skill bars reveal
var sk=document.querySelector('.skills');
if('IntersectionObserver' in window){
  new IntersectionObserver(function(en,o){if(en[0].isIntersecting){sk.classList.add('vis');o.disconnect()}},{threshold:.2}).observe(sk);
}else sk.classList.add('vis');

// Arranque: escribe ./AlexDBdevopsCV.sh, espera Enter (o arranca solo) y va mostrando el CV
(function(){
  var html=document.documentElement;
  if(!html.classList.contains('boot'))return;
  if(!document.body.classList.contains('v-a')){html.classList.remove('boot');return}
  var cmd=document.querySelector('.bootcmd'),log=document.querySelector('.bootlog'),btn=document.querySelector('.bootenter'),
      items=[].slice.call(document.querySelectorAll('.bt')),typed=[].slice.call(document.querySelectorAll('.cmdline .typed')),
      full=cmd.textContent,t0,timers=[],state='typing',autoTimer;
  cmd.textContent='';typed.forEach(function(s){s.dataset.full=s.textContent;s.textContent=''});
  function at(ms,fn){timers.push(setTimeout(fn,ms))}
  function show(sel){items.filter(function(e){return e.matches(sel)}).forEach(function(e){e.classList.add('in')})}
  function finish(){
    if(state==='done')return;state='done';timers.forEach(clearTimeout);clearTimeout(autoTimer);btn.hidden=true;
    cmd.textContent=full;typed.forEach(function(s){s.textContent=s.dataset.full});
    items.forEach(function(e){e.classList.add('in')});
    log.textContent='finish in '+((performance.now()-(t0||performance.now()))/1000).toFixed(1)+'s';
    setTimeout(function(){html.classList.remove('boot');log.textContent=''},2200);
    try{sessionStorage.setItem('booted','1')}catch(e){}
  }
  function run(){
    if(state!=='waiting')return;state='running';clearTimeout(autoTimer);btn.hidden=true;t0=performance.now();
    var steps=[[0,'.stage','rendering eye.svg'],[350,'h1','loading profile'],[480,'.role',''],[650,'.cmdline',''],
      [1250,'.banner','echo'],[1550,'.intro',''],[1700,'.pills, .lang','building sections'],
      [1950,'.experience','experience (es/en)'],[2250,'.skills','skills'],[2550,'.build','how it was built'],[2750,'.contact, footer','']];
    steps.forEach(function(s){at(s[0],function(){show(s[1]);if(s[2])log.textContent='> '+s[2]+'…'})});
    typed.forEach(function(s){var f=s.dataset.full;for(var k=1;k<=f.length;k++){(function(n){at(700+n*32,function(){s.textContent=f.slice(0,n)})})(k)}});
    at(3300,finish);
  }
  addEventListener('keydown',function(e){
    if(state==='waiting'&&(e.key==='Enter'||e.key===' ')){e.preventDefault();run()}
    else if(state==='running'||state==='typing'){finish()}
  });
  btn.addEventListener('click',function(e){e.stopPropagation();run()});
  addEventListener('pointerdown',function(e){if(e.target!==btn&&(state==='running'||state==='typing'))finish()});
  var t=350;
  for(var i=1;i<=full.length;i++){(function(n){at(t+n*55,function(){cmd.textContent=full.slice(0,n)})})(i)}
  at(t+full.length*55+200,function(){
    state='waiting';btn.hidden=false;
    autoTimer=setTimeout(run,3200);   // si nadie pulsa Enter, se genera solo
  });
})();

// Copiar el correo
[].slice.call(document.querySelectorAll('.copy')).forEach(function(b){
  var label=b.textContent;
  b.addEventListener('click',function(){
    var done=function(){b.textContent=root.lang==='es'?'copiado':'copied';setTimeout(function(){b.textContent=label},1500)};
    if(navigator.clipboard)navigator.clipboard.writeText(b.dataset.copy).then(done,function(){});
  });
});
})();

// A reload (or coming back to the tab from the browser cache) starts at the top, at the prompt,
// instead of where the page was left; links to a section (#...) keep their position
(function(){
  function top(){ if(!location.hash) window.scrollTo(0,0) }
  top();
  window.addEventListener('load',top);
  window.addEventListener('pageshow',function(e){ if(e.persisted) top() });
})();
