const state = { mailbox: "", messages: [], selectedUID: null };

function trunc(s, n){
  if(!s) return "";
  s = s.replace(/\s+/g,' ').trim();
  return s.length>n ? s.slice(0,n-1)+'…' : s;
}

function stripHtml(html){
  if(!html) return '';
  const tmp = document.createElement('div');
  tmp.innerHTML = html;
  const text = (tmp.textContent || tmp.innerText || '').replace(/\s+/g,' ').trim();
  return text;
}

function fmtTime(ts){
  if(!ts) return "";
  const d = new Date(ts*1000);
  return d.toLocaleString();
}

async function fetchMessages(){
  if(state._fetching) return; state._fetching = true;
  try{
    const res = await fetch('/messages',{cache:'no-store'});
    if(!res.ok) throw new Error('failed');
    const data = await res.json();
    state.mailbox = data.mailbox || '';
    // Preserve already-fetched bodies to avoid reloading heavy HTML/images
  const prev = Object.create(null);
  for(const x of state.messages){ prev[x.uid] = x; }
  const fresh = (data.data || []).map(n => {
    const old = prev[n.uid];
    if(old && old.bodyHtml){ n.bodyHtml = old.bodyHtml; }
    if(old && old.attachments){ n.attachments = old.attachments; }
    return n;
  });
    state.messages = fresh;
    // signature to detect meaningful changes; skip DOM work if same
    const sig = fresh.map(m=> (m.uid||'') + ':' + (m.receivedAt||0)).join('|');
    const changed = sig !== state._sig;
    state._sig = sig;
    document.getElementById('current-email').textContent = state.mailbox || 'unknown';
    const inp = document.getElementById('current-email-input');
    if(inp) inp.value = state.mailbox || '';
    const cnt = document.getElementById('inbox-count');
    if(cnt) cnt.textContent = state.messages.length ? String(state.messages.length) : '';
    if(changed) renderList();
    if(state.selectedUID==null && state.messages.length){
      selectMessage(state.messages[0].uid);
    }
  }finally{
    state._fetching = false;
  }
}

function renderList(){
  const ul = document.getElementById('messages');
  ul.innerHTML = '';
  const empty = document.getElementById('messages-empty');
  if(empty){
    const noMail = !state.messages.length;
    // Ensure visibility toggles correctly regardless of CSS precedence
    empty.hidden = !noMail ? true : false;     // hide when có thư
    empty.style.display = noMail ? 'flex' : 'none';
  }
  if(!state.messages.length){ return; }
  for(const m of state.messages){
    const li = document.createElement('li');
    li.dataset.uid = m.uid;
    if(m.uid===state.selectedUID) li.classList.add('active');

    const wrap = document.createElement('div');
    wrap.className = 'message';

    const av = document.createElement('div');
    av.className = 'avatar';
    const from = (m.from||'').trim();
    const letter = (from.match(/[a-zA-Z0-9]/)||['?'])[0].toUpperCase();
    av.textContent = letter;

    const info = document.createElement('div');
    const row = document.createElement('div'); row.className='sender-time';
    const sender = document.createElement('div'); sender.className='sender'; sender.textContent = trunc(from, 28);
    const time = document.createElement('div'); time.className='time'; time.textContent = fmtTime(m.receivedAt);
    row.appendChild(sender); row.appendChild(time);

    const subj = document.createElement('div'); subj.className='subject'; subj.textContent = trunc(m.subject||'', 80);
    const snippet = document.createElement('div'); snippet.className='snippet';
    const plain = (m.bodyPreview || stripHtml(m.bodyHtml||''));
    snippet.textContent = trunc(plain || '', 120);

    info.appendChild(row); info.appendChild(subj); info.appendChild(snippet);
    wrap.appendChild(av); wrap.appendChild(info);
    li.appendChild(wrap);
    li.addEventListener('click',()=> selectMessage(m.uid));
    ul.appendChild(li);
  }
}

function selectMessage(uid){
  state.selectedUID = uid;
  const m = state.messages.find(x=>x.uid===uid);
  if(!m) return;
  document.getElementById('detail-subject').textContent = m.subject || '';
  document.getElementById('detail-from').textContent = m.from || '';
  document.getElementById('detail-date').textContent = fmtTime(m.receivedAt);
  const frame = document.getElementById('body-frame');
  // Reset frame height so không giữ chiều cao của thư trước
  try{ frame.style.height = '360px'; }catch(_){ }

  // helper: render attachments block from message
  const renderAttachmentsBlock = (msg) => {
    const att = document.getElementById('attachments');
    if(!att) return;
    att.innerHTML = '';
    if(Array.isArray(msg.attachments) && msg.attachments.length){
      const fmtBytes = (b)=>{ if(!b && b!==0) return ''; const u=['B','KB','MB','GB']; let i=0, n=b; while(n>=1024 && i<u.length-1){ n/=1024; i++; } return ` (${n.toFixed(n>=10?0:1)} ${u[i]})`; };
      for(const a of msg.attachments){
        const id = a._id||a.id||a.ID||a.Id||a.Id;
        if(!id) continue;
        const ael = document.createElement('a');
        ael.href = `/attachments/${id}`;
        ael.target = '_blank';
        ael.setAttribute('download', a.filename || 'attachment');
        const name = a.filename || a.Filename || a.cid || 'attachment';
        ael.textContent = name + fmtBytes(a.size||a.Size);
        ael.title = `Tải ${name}`;
        att.appendChild(ael);
      }
    }
  };
  // first pass: render any known attachments (from detail cache if available)
  renderAttachmentsBlock(m);

  // Ensure attachments exist: nếu danh sách rỗng mà đếm > 0, fetch chi tiết để cập nhật
  if((!m.attachments || m.attachments.length===0) && (m.attachmentsCount>0 || m.AttcCount>0)){
    if(!m._attLoading){
      m._attLoading = true;
      (async()=>{
        try{
          const res = await fetch(`/messages/${m.uid}`, {cache:'no-store'});
          if(!res.ok) return;
          const full = await res.json();
          if(state.selectedUID !== m.uid) return;
          if(full.attachments){ m.attachments = full.attachments; renderAttachmentsBlock(m); }
        }catch(_){ } finally { m._attLoading = false; }
      })();
    }
  }
  const setFrame = (html) => {
    const nonce = Date.now().toString(36) + Math.random().toString(36).slice(2);
    frame.srcdoc = `<!DOCTYPE html><html><head><meta charset='utf-8'>
    <meta name='viewport' content='width=device-width, initial-scale=1'>
    <meta name='x-refresh' content='${nonce}'>
    <style>
      html,body{margin:0;padding:12px;font-family:-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;font-size:16px;line-height:1.6;color:#111;}
      body{-webkit-text-size-adjust:100%;text-size-adjust:100%;overflow-wrap:anywhere;word-break:break-word}
      img,video{max-width:100%;height:auto}
      table{max-width:100%;border-collapse:collapse}
      .table-wrap{overflow:auto}
      a{color:#2563eb}
      @media (min-width: 1024px){ body{ font-size:17px } }
      @media (min-width: 1440px){ body{ font-size:18px } }
    </style>
    </head><body>${html||''}</body></html>`;
  };
  if(m.bodyHtml){
    setFrame(m.bodyHtml);
    // đo ngay sau khi set để tránh phải đợi onload
    setTimeout(()=>{
      try{
        const doc = frame.contentDocument || frame.contentWindow.document; if(!doc) return;
        frame.style.height = Math.max(doc.body?.scrollHeight||0, doc.documentElement?.scrollHeight||0, 320) + 'px';
      }catch(_){ }
    }, 0);
    state._loadedUID = uid;
  }else{
    // lightweight placeholder, then fetch detail only once
    setFrame('<p style="color:#64748b">Đang tải nội dung…</p>');
    (async ()=>{
      try{
        // Hủy request trước đó nếu user đổi thư rất nhanh
        try{ state._detailAbort && state._detailAbort.abort && state._detailAbort.abort(); }catch(_){ }
        const ctl = new AbortController(); state._detailAbort = ctl;
        const res = await fetch(`/messages/${m.uid}`, {cache:'no-store', signal: ctl.signal});
        if(!res.ok) return;
        const full = await res.json();
        // ensure still same selection
        if(state.selectedUID !== m.uid) return;
        m.bodyHtml = full.bodyHtml || full.Body || '';
        // merge attachments if provided
        if(full.attachments){ m.attachments = full.attachments; renderAttachmentsBlock(m); }
        setFrame(m.bodyHtml);
        setTimeout(()=>{
          try{
            const doc = frame.contentDocument || frame.contentWindow.document; if(!doc) return;
            frame.style.height = Math.max(doc.body?.scrollHeight||0, doc.documentElement?.scrollHeight||0, 320) + 'px';
          }catch(_){ }
        }, 0);
        state._loadedUID = uid;
      }catch(_){ }
    })();
  }
  const adjustFrame = () => {
    try{
      const doc = frame.contentDocument || frame.contentWindow.document;
      if(!doc) return;
      const body = doc.body, htmlEl = doc.documentElement;
      const h = Math.max(
        body ? body.scrollHeight : 0,
        htmlEl ? htmlEl.scrollHeight : 0,
        320
      );
      if(frame._lastHeight !== h){ frame.style.height = h + 'px'; frame._lastHeight = h; }
    }catch(e){}
  };
  frame.onload = () => {
    adjustFrame();
    try{
      const doc = frame.contentDocument || frame.contentWindow.document;
      // Hint browsers to avoid aggressive network usage for images
      try{
        const imgs = doc && doc.querySelectorAll ? doc.querySelectorAll('img') : [];
        imgs && imgs.forEach && imgs.forEach(img => {
          try{
            if(!img.hasAttribute('loading')) img.setAttribute('loading','lazy');
            if(!img.hasAttribute('decoding')) img.setAttribute('decoding','async');
            if(!img.hasAttribute('referrerpolicy')) img.setAttribute('referrerpolicy','no-referrer');
            if(!img.hasAttribute('fetchpriority')) img.setAttribute('fetchpriority','low');
          }catch(_){ }
        });
      }catch(_){ }
      // cleanup previous observers if any
      try{ frame._mo && frame._mo.disconnect && frame._mo.disconnect(); }catch(_){ }
      try{ frame._ro && frame._ro.disconnect && frame._ro.disconnect(); }catch(_){ }
      // Mutation observer for DOM changes
      try{
        const mo = new MutationObserver(()=> adjustFrame());
        mo.observe(doc.documentElement, {subtree:true, childList:true, attributes:true, characterData:true});
        frame._mo = mo;
      }catch(_){ }
      // ResizeObserver for layout changes
      try{
        const RO = doc.defaultView && doc.defaultView.ResizeObserver;
        if(RO){
          const ro = new RO(()=> adjustFrame());
          if(doc.documentElement) ro.observe(doc.documentElement);
          if(doc.body) ro.observe(doc.body);
          frame._ro = ro;
        }
      }catch(_){ }
    }catch(_){ }
    // multiple passes for late assets
    setTimeout(adjustFrame, 150);
    setTimeout(adjustFrame, 400);
    setTimeout(adjustFrame, 800);
    setTimeout(adjustFrame, 1600);
    // rAF loop vài nhịp để bắt các layout pass
    let n=0; const raf = () => { if(n++<5){ requestAnimationFrame(()=>{ adjustFrame(); raf(); }); } }; raf();
  };

  const prev = document.getElementById('preview-text');
  // hide separate preview block to avoid big gaps
  prev.style.display = 'none';

  // attachments already rendered above
  // do not re-render list here to avoid churn
}

async function randomize(){
  const res = await fetch('/randomize',{method:'POST'});
  if(res.ok){
    await fetchMessages();
  }
}

function applyTheme(theme){
  if(theme === 'dark'){
    document.documentElement.setAttribute('data-theme','dark');
  }else{
    document.documentElement.removeAttribute('data-theme');
  }
}

function initTheme(){
  let theme = localStorage.getItem('theme');
  if(!theme){
    const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
    theme = prefersDark ? 'dark' : 'light';
  }
  applyTheme(theme);
  const btn = document.getElementById('theme-toggle');
  if(btn){
    const label = () => (document.documentElement.getAttribute('data-theme')==='dark' ? 'Sáng' : 'Tối');
    const update = () => { const el = btn.querySelector('.btn-label'); if(el) el.textContent = label(); };
    update();
    btn.addEventListener('click', ()=>{
      const isDark = document.documentElement.getAttribute('data-theme')==='dark';
      const next = isDark ? 'light' : 'dark';
      applyTheme(next);
      localStorage.setItem('theme', next);
      update();
    });
  }
}

async function main(){
  // show spinner while refreshing
  const refreshBtn = document.getElementById('refresh-btn');
  async function fetchWithSpinner(){
    try{ if(refreshBtn) refreshBtn.setAttribute('data-loading','1'); await fetchMessages(); }
    finally { if(refreshBtn) refreshBtn.removeAttribute('data-loading'); }
  }
  if(refreshBtn){ refreshBtn.addEventListener('click', fetchWithSpinner); }
  document.getElementById('random-btn').addEventListener('click', randomize);
  document.getElementById('copy-btn').addEventListener('click', async ()=>{
    if(!state.mailbox) return;
    let ok=false;
    try{
      await navigator.clipboard.writeText(state.mailbox);
      ok=true;
    }catch(e){
      // fallback
      const ta = document.createElement('textarea');
      ta.value = state.mailbox; document.body.appendChild(ta); ta.select();
      try{ ok = document.execCommand('copy'); }catch(_){}
      document.body.removeChild(ta);
    }
    const btn = document.getElementById('copy-btn');
    const lbl = btn ? btn.querySelector('.btn-label') : null;
    if(lbl){
      const old = lbl.textContent;
      lbl.textContent = ok? 'Đã sao chép!' : 'Sao chép thất bại';
      setTimeout(()=> lbl.textContent = old, 1200);
    }
  });
  try{ await fetchMessages(); }catch(e){ console.error(e); }
  // Poll slower and pause when tab hidden to avoid waste
  setInterval(()=>{ if(document.visibilityState !== 'hidden') fetchMessages(); }, 5000);
}

window.addEventListener('DOMContentLoaded', ()=>{ initTheme(); main(); });
