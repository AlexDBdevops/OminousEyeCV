(function(){
var root=document.documentElement,svg=document.getElementById('art'),iris=document.getElementById('iris'),
disc=document.getElementById('disc'),sectors=[].slice.call(document.querySelectorAll('.sector')),
pills=[].slice.call(document.querySelectorAll('.pill')),details=[].slice.call(document.querySelectorAll('.detail'));
var reduce=matchMedia('(prefers-reduced-motion: reduce)').matches;

// Eye tracking
function look(x,y){
  var r=svg.getBoundingClientRect(),cx=r.left+r.width/2,cy=r.top+r.height/2;
  var dx=x-cx,dy=y-cy,a=Math.atan2(dy,dx),d=Math.min(Math.hypot(dx,dy)/(r.width*0.5),1);
  var k=r.width/600,m=34*d;
  iris.style.transform='translate('+(Math.cos(a)*m)+'px,'+(Math.sin(a)*m*0.55)+'px)';
}
if(!reduce){
  addEventListener('pointermove',function(e){look(e.clientX,e.clientY)},{passive:true});
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
})();
