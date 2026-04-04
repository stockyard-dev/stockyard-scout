package server
import "net/http"
func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "text/html"); w.Write([]byte(dashHTML)) }
const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Scout</title><link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet"><style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}.main{padding:1.5rem;max-width:960px;margin:0 auto}.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}.toolbar{display:flex;gap:.5rem;margin-bottom:1rem}.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.search:focus{outline:none;border-color:var(--leather)}.target{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}.target:hover{border-color:var(--leather)}.target-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}.target-name{font-size:.85rem;font-weight:700}.target-url{font-size:.6rem;color:var(--blue);margin-top:.1rem;word-break:break-all}.target-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}.target-actions{display:flex;gap:.3rem;flex-shrink:0}.dot{display:inline-block;width:8px;height:8px;border-radius:50%;margin-right:.3rem}.dot.up{background:var(--green)}.dot.down{background:var(--red)}.dot.unknown{background:var(--cm)}.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}.badge.up{border-color:var(--green);color:var(--green)}.badge.down{border-color:var(--red);color:var(--red)}.badge.paused{border-color:var(--cm);color:var(--cm)}.latency{font-size:.6rem}.latency.good{color:var(--green)}.latency.warn{color:var(--gold)}.latency.bad{color:var(--red)}.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}.btn-sm{font-size:.55rem;padding:.2rem .4rem}.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw}.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}.fr input,.fr select{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.fr input:focus,.fr select:focus{outline:none;border-color:var(--leather)}.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> SCOUT</h1><button class="btn btn-p" onclick="openForm()">+ Add Target</button></div>
<div class="main"><div class="stats" id="stats"></div><div class="toolbar"><input class="search" id="search" placeholder="Search targets..." oninput="render()"></div><div id="list"></div></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/targets').then(function(r){return r.json()});items=r.targets||[];renderStats();render();}
function renderStats(){var t=items.length,up=items.filter(function(t){return t.status==='up'}).length,down=items.filter(function(t){return t.status==='down'}).length;
document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+t+'</div><div class="st-l">Targets</div></div><div class="st"><div class="st-v" style="color:var(--green)">'+up+'</div><div class="st-l">Up</div></div><div class="st"><div class="st-v" style="color:'+(down>0?'var(--red)':'var(--cream)')+'">'+down+'</div><div class="st-l">Down</div></div>';}
function render(){var q=(document.getElementById('search').value||'').toLowerCase();var f=items;
if(q)f=f.filter(function(t){return(t.name||'').toLowerCase().includes(q)||(t.url||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No monitoring targets.</div>';return;}
f.sort(function(a,b){return a.status==='down'?-1:b.status==='down'?1:0});
var h='';f.forEach(function(t){
var st=t.status||'unknown';
h+='<div class="target"><div class="target-top"><div style="flex:1">';
h+='<div class="target-name"><span class="dot '+st+'"></span>'+esc(t.name)+'</div>';
h+='<div class="target-url">'+esc(t.url)+'</div>';
h+='</div><div class="target-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(''+t.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+t.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div><div class="target-meta">';
h+='<span class="badge '+st+'">'+st+'</span>';
if(t.type)h+='<span>'+esc(t.type)+'</span>';
if(t.interval_seconds)h+='<span>every '+t.interval_seconds+'s</span>';
if(t.last_latency_ms!=null){var cls=t.last_latency_ms<200?'good':t.last_latency_ms<1000?'warn':'bad';h+='<span class="latency '+cls+'">'+t.last_latency_ms+'ms</span>';}
if(t.last_checked_at)h+='<span>checked: '+ft(t.last_checked_at)+'</span>';
h+='</div></div>';});
document.getElementById('list').innerHTML=h;}
async function del(id){if(!confirm('Remove?'))return;await fetch(A+'/targets/'+id,{method:'DELETE'});load();}
function formHTML(tgt){var i=tgt||{name:'',url:'',type:'http',interval_seconds:60};var isEdit=!!tgt;
var h='<h2>'+(isEdit?'EDIT':'ADD')+' TARGET</h2>';
h+='<div class="fr"><label>Name *</label><input id="f-name" value="'+esc(i.name)+'" placeholder="My API"></div>';
h+='<div class="fr"><label>URL *</label><input id="f-url" value="'+esc(i.url)+'" placeholder="https://api.example.com/health"></div>';
h+='<div class="row2"><div class="fr"><label>Type</label><select id="f-type"><option value="http"'+(i.type==='http'?' selected':'')+'>HTTP</option><option value="tcp"'+(i.type==='tcp'?' selected':'')+'>TCP</option><option value="ping"'+(i.type==='ping'?' selected':'')+'>Ping</option></select></div>';
h+='<div class="fr"><label>Interval (sec)</label><input id="f-int" type="number" value="'+(i.interval_seconds||60)+'"></div></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button></div>';return h;}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var t=null;for(var j=0;j<items.length;j++){if(items[j].id===id){t=items[j];break;}}if(!t)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(t);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var name=document.getElementById('f-name').value.trim();var url=document.getElementById('f-url').value.trim();if(!name||!url){alert('Name and URL required');return;}
var body={name:name,url:url,type:document.getElementById('f-type').value,interval_seconds:parseInt(document.getElementById('f-int').value)||60};
if(editId){await fetch(A+'/targets/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/targets',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}closeModal();load();}
function ft(t){if(!t)return'';try{var d=new Date(t);return d.toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});load();
</script></body></html>`
