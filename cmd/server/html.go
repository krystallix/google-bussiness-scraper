package main

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>G-Lead Scraper</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Geist+Mono:wght@400;500;600&family=Geist:wght@400;500;600;700&display=swap" rel="stylesheet">
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

  :root {
    --bg:       #09090b;
    --surface:  #18181b;
    --border:   #27272a;
    --border2:  #3f3f46;
    --text:     #fafafa;
    --muted:    #a1a1aa;
    --accent:   #22c55e;
    --accent-d: #16a34a;
    --amber:    #f59e0b;
    --red:      #ef4444;
    --blue:     #3b82f6;
    --sidebar:  220px;
    --radius:   8px;
  }

  html, body { height: 100%; background: var(--bg); color: var(--text); font-family: 'Geist', sans-serif; font-size: 14px; }

  /* ── Layout ── */
  .layout { display: flex; height: 100vh; overflow: hidden; }

  /* ── Sidebar ── */
  .sidebar {
    width: var(--sidebar); flex-shrink: 0;
    background: var(--surface); border-right: 1px solid var(--border);
    display: flex; flex-direction: column; padding: 20px 0;
  }
  .logo {
    display: flex; align-items: center; gap: 10px;
    padding: 0 20px 24px; border-bottom: 1px solid var(--border);
    margin-bottom: 12px;
  }
  .logo-icon {
    width: 32px; height: 32px; background: var(--accent);
    border-radius: 6px; display: flex; align-items: center; justify-content: center;
    font-size: 16px; flex-shrink: 0;
  }
  .logo-text { font-weight: 700; font-size: 15px; letter-spacing: -0.02em; line-height: 1.1; }
  .logo-text span { display: block; font-size: 10px; color: var(--muted); font-weight: 400; letter-spacing: 0.06em; text-transform: uppercase; font-family: 'Geist Mono', monospace; }

  .nav { display: flex; flex-direction: column; gap: 2px; padding: 0 10px; }
  .nav-item {
    display: flex; align-items: center; gap: 10px;
    padding: 9px 10px; border-radius: var(--radius);
    color: var(--muted); cursor: pointer; transition: background 0.15s, color 0.15s;
    font-size: 13px; font-weight: 500; user-select: none;
  }
  .nav-item:hover { background: var(--border); color: var(--text); }
  .nav-item.active { background: rgba(34,197,94,0.12); color: var(--accent); }
  .nav-item svg { width: 16px; height: 16px; flex-shrink: 0; }

  .sidebar-footer { margin-top: auto; padding: 16px 20px 0; border-top: 1px solid var(--border); }
  .sidebar-footer p { font-size: 11px; color: var(--muted); font-family: 'Geist Mono', monospace; }

  /* ── Main ── */
  .main { flex: 1; display: flex; flex-direction: column; overflow: hidden; }

  /* ── Topbar ── */
  .topbar {
    padding: 16px 24px; border-bottom: 1px solid var(--border);
    background: var(--surface);
    display: flex; align-items: center; gap: 12px; flex-wrap: wrap;
  }
  .topbar-form { display: flex; align-items: center; gap: 10px; flex: 1; }
  .input-wrap { position: relative; flex: 1; max-width: 520px; }
  .input-wrap svg { position: absolute; left: 12px; top: 50%; transform: translateY(-50%); width: 15px; height: 15px; color: var(--muted); pointer-events: none; }
  .keyword-input {
    width: 100%; padding: 10px 12px 10px 36px;
    background: var(--bg); border: 1px solid var(--border2);
    border-radius: var(--radius); color: var(--text); font-family: 'Geist', sans-serif;
    font-size: 14px; outline: none; transition: border-color 0.15s;
  }
  .keyword-input:focus { border-color: var(--accent); }
  .keyword-input::placeholder { color: var(--muted); }

  .limit-wrap { display: flex; align-items: center; gap: 8px; }
  .limit-wrap label { font-size: 12px; color: var(--muted); white-space: nowrap; font-family: 'Geist Mono', monospace; }
  .limit-input {
    width: 72px; padding: 10px 10px;
    background: var(--bg); border: 1px solid var(--border2);
    border-radius: var(--radius); color: var(--text); font-family: 'Geist Mono', monospace;
    font-size: 14px; outline: none; text-align: center; transition: border-color 0.15s;
  }
  .limit-input:focus { border-color: var(--accent); }

  .headful-wrap { display: flex; align-items: center; gap: 6px; }
  .headful-wrap label { font-size: 12px; color: var(--muted); white-space: nowrap; }
  .toggle { position: relative; width: 36px; height: 20px; cursor: pointer; }
  .toggle input { opacity: 0; width: 0; height: 0; }
  .toggle-track {
    position: absolute; inset: 0; background: var(--border2);
    border-radius: 20px; transition: background 0.2s;
  }
  .toggle input:checked + .toggle-track { background: var(--accent); }
  .toggle-thumb {
    position: absolute; top: 3px; left: 3px; width: 14px; height: 14px;
    background: #fff; border-radius: 50%; transition: transform 0.2s;
  }
  .toggle input:checked ~ .toggle-thumb { transform: translateX(16px); }

  .btn-scrape {
    padding: 10px 20px; background: var(--accent); color: #000;
    border: none; border-radius: var(--radius); font-family: 'Geist', sans-serif;
    font-size: 13px; font-weight: 700; cursor: pointer; white-space: nowrap;
    transition: background 0.15s, transform 0.1s; letter-spacing: 0.01em;
  }
  .btn-scrape:hover { background: var(--accent-d); }
  .btn-scrape:active { transform: scale(0.98); }
  .btn-scrape:disabled { opacity: 0.5; cursor: not-allowed; transform: none; }

  /* ── Status bar ── */
  .statusbar {
    padding: 8px 24px; background: var(--bg); border-bottom: 1px solid var(--border);
    display: flex; align-items: center; gap: 16px;
    font-size: 12px; font-family: 'Geist Mono', monospace; color: var(--muted);
    min-height: 36px;
  }
  .status-dot {
    width: 7px; height: 7px; border-radius: 50%; background: var(--border2); flex-shrink: 0;
  }
  .status-dot.running { background: var(--amber); animation: pulse 1s infinite; }
  .status-dot.done { background: var(--accent); }
  .status-dot.error { background: var(--red); }
  @keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.4; } }

  /* ── Content ── */
  .content { flex: 1; overflow: auto; padding: 20px 24px; }

  /* ── Filter bar ── */
  .filter-bar { display: flex; align-items: center; gap: 8px; margin-bottom: 16px; flex-wrap: wrap; }
  .filter-chip {
    padding: 5px 14px; border-radius: 999px; font-size: 12px; font-weight: 500;
    border: 1px solid var(--border2); background: transparent; color: var(--muted);
    cursor: pointer; transition: all 0.15s; font-family: 'Geist Mono', monospace;
  }
  .filter-chip:hover { border-color: var(--text); color: var(--text); }
  .filter-chip.active { background: var(--text); color: var(--bg); border-color: var(--text); }
  .filter-chip.chip-high.active { background: var(--accent); color: #000; border-color: var(--accent); }
  .filter-chip.chip-standard.active { background: var(--blue); color: #fff; border-color: var(--blue); }
  .filter-chip.chip-skip.active { background: var(--red); color: #fff; border-color: var(--red); }

  .filter-right { margin-left: auto; display: flex; align-items: center; gap: 8px; }
  .count-badge {
    font-size: 11px; font-family: 'Geist Mono', monospace; color: var(--muted);
    background: var(--surface); padding: 3px 10px; border-radius: 999px;
    border: 1px solid var(--border);
  }
  .search-input {
    padding: 6px 12px; background: var(--surface); border: 1px solid var(--border);
    border-radius: var(--radius); color: var(--text); font-family: 'Geist', sans-serif;
    font-size: 13px; outline: none; width: 200px; transition: border-color 0.15s;
  }
  .search-input:focus { border-color: var(--accent); }
  .search-input::placeholder { color: var(--muted); }

  /* ── Table ── */
  .table-wrap { border: 1px solid var(--border); border-radius: var(--radius); overflow: hidden; }
  table { width: 100%; border-collapse: collapse; }
  thead { background: var(--surface); border-bottom: 1px solid var(--border); }
  th {
    padding: 10px 14px; text-align: left; font-size: 11px; font-weight: 600;
    color: var(--muted); letter-spacing: 0.06em; text-transform: uppercase;
    font-family: 'Geist Mono', monospace; white-space: nowrap; cursor: pointer;
    user-select: none;
  }
  th:hover { color: var(--text); }
  th .sort-icon { display: inline; opacity: 0.4; margin-left: 4px; }
  th.sorted .sort-icon { opacity: 1; color: var(--accent); }

  tbody tr {
    border-bottom: 1px solid var(--border);
    transition: background 0.1s;
  }
  tbody tr:last-child { border-bottom: none; }
  tbody tr:hover { background: rgba(255,255,255,0.03); }
  td { padding: 11px 14px; font-size: 13px; vertical-align: middle; }

  .cell-name { font-weight: 500; max-width: 220px; }
  .cell-name a { color: var(--text); text-decoration: none; }
  .cell-name a:hover { color: var(--accent); }
  .cell-name .sub { font-size: 11px; color: var(--muted); margin-top: 2px; font-family: 'Geist Mono', monospace; }

  .stars { color: var(--amber); font-size: 12px; letter-spacing: -1px; }
  .rating-num { font-size: 12px; color: var(--muted); font-family: 'Geist Mono', monospace; margin-left: 4px; }

  .reviews-num { font-family: 'Geist Mono', monospace; font-size: 13px; }

  .chip {
    display: inline-flex; align-items: center; gap: 4px;
    padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 500;
    font-family: 'Geist Mono', monospace; white-space: nowrap;
  }
  .chip-nosite { background: rgba(245,158,11,0.15); color: var(--amber); border: 1px solid rgba(245,158,11,0.3); }
  .chip-hassite { background: rgba(63,63,70,0.6); color: var(--muted); border: 1px solid var(--border2); }

  .tier-badge {
    display: inline-block; padding: 3px 10px; border-radius: 4px;
    font-size: 11px; font-weight: 700; font-family: 'Geist Mono', monospace;
    letter-spacing: 0.04em; white-space: nowrap;
  }
  .tier-HIGH { background: rgba(34,197,94,0.15); color: var(--accent); border: 1px solid rgba(34,197,94,0.3); }
  .tier-STANDARD { background: rgba(59,130,246,0.15); color: var(--blue); border: 1px solid rgba(59,130,246,0.3); }
  .tier-SKIP { background: rgba(239,68,68,0.1); color: #f87171; border: 1px solid rgba(239,68,68,0.2); opacity: 0.7; }

  .btn-wa {
    display: inline-flex; align-items: center; gap: 6px;
    padding: 5px 10px; border-radius: 6px; font-size: 11px; font-weight: 600;
    background: rgba(37,211,102,0.12); color: #25d366; border: 1px solid rgba(37,211,102,0.3);
    text-decoration: none; cursor: pointer; transition: background 0.15s;
    font-family: 'Geist Mono', monospace; white-space: nowrap;
  }
  .btn-wa:hover { background: rgba(37,211,102,0.22); }
  .btn-wa svg { width: 12px; height: 12px; }

  /* ── Empty / Loading ── */
  .empty-state {
    display: flex; flex-direction: column; align-items: center; justify-content: center;
    padding: 80px 20px; color: var(--muted); gap: 12px; text-align: center;
  }
  .empty-state svg { width: 40px; height: 40px; opacity: 0.3; }
  .empty-state p { font-size: 14px; }
  .empty-state small { font-size: 12px; opacity: 0.6; font-family: 'Geist Mono', monospace; }

  .loading-bar {
    height: 2px; background: var(--border);
    position: relative; overflow: hidden; border-radius: 1px; margin-bottom: 16px;
  }
  .loading-bar::after {
    content: ''; position: absolute; top: 0; left: -60%;
    width: 60%; height: 100%; background: var(--accent);
    animation: slide 1.2s ease-in-out infinite;
  }
  @keyframes slide { to { left: 100%; } }

  /* ── WhatsApp Modal ── */
  .modal-overlay {
    display: none; position: fixed; inset: 0; background: rgba(0,0,0,0.7);
    z-index: 100; align-items: center; justify-content: center;
    backdrop-filter: blur(4px);
  }
  .modal-overlay.open { display: flex; }
  .modal {
    background: var(--surface); border: 1px solid var(--border2);
    border-radius: 12px; padding: 24px; width: 90%; max-width: 520px;
    display: flex; flex-direction: column; gap: 16px;
  }
  .modal-header { display: flex; align-items: center; justify-content: space-between; }
  .modal-title { font-size: 15px; font-weight: 600; }
  .modal-close { background: none; border: none; color: var(--muted); cursor: pointer; padding: 4px; font-size: 18px; line-height: 1; }
  .modal-close:hover { color: var(--text); }
  .modal label { font-size: 12px; color: var(--muted); display: block; margin-bottom: 6px; font-family: 'Geist Mono', monospace; }
  .modal textarea {
    width: 100%; padding: 10px 12px; background: var(--bg); border: 1px solid var(--border2);
    border-radius: var(--radius); color: var(--text); font-family: 'Geist', sans-serif;
    font-size: 13px; outline: none; resize: vertical; min-height: 120px;
    transition: border-color 0.15s; line-height: 1.5;
  }
  .modal textarea:focus { border-color: var(--accent); }
  .modal-footer { display: flex; gap: 8px; justify-content: flex-end; }
  .btn-ghost {
    padding: 8px 16px; background: transparent; border: 1px solid var(--border2);
    border-radius: var(--radius); color: var(--muted); cursor: pointer;
    font-family: 'Geist', sans-serif; font-size: 13px; transition: all 0.15s;
  }
  .btn-ghost:hover { border-color: var(--text); color: var(--text); }
  .btn-send {
    padding: 8px 16px; background: #25d366; color: #000;
    border: none; border-radius: var(--radius); cursor: pointer;
    font-family: 'Geist', sans-serif; font-size: 13px; font-weight: 700; transition: opacity 0.15s;
  }
  .btn-send:hover { opacity: 0.85; }

  /* ── Tooltip ── */
  [data-tip] { position: relative; cursor: help; }
  [data-tip]::after {
    content: attr(data-tip); position: absolute; bottom: 100%; left: 50%; transform: translateX(-50%);
    background: #000; color: #fff; font-size: 11px; padding: 4px 8px; border-radius: 4px;
    white-space: nowrap; pointer-events: none; opacity: 0; transition: opacity 0.15s;
    margin-bottom: 6px; font-family: 'Geist Mono', monospace; z-index: 10;
  }
  [data-tip]:hover::after { opacity: 1; }

  /* Scrollbar */
  ::-webkit-scrollbar { width: 6px; height: 6px; }
  ::-webkit-scrollbar-track { background: transparent; }
  ::-webkit-scrollbar-thumb { background: var(--border2); border-radius: 3px; }
</style>
</head>
<body>

<div class="layout">
  <!-- Sidebar -->
  <aside class="sidebar">
    <div class="logo">
      <div class="logo-icon">
        <svg viewBox="0 0 20 20" fill="none" stroke="#000" stroke-width="2" style="width:18px;height:18px">
          <circle cx="10" cy="10" r="3"/><path d="M10 2v2M10 16v2M2 10h2M16 10h2"/>
          <path d="M4.93 4.93l1.41 1.41M13.66 13.66l1.41 1.41M4.93 15.07l1.41-1.41M13.66 6.34l1.41-1.41"/>
        </svg>
      </div>
      <div class="logo-text">G-Lead<span>Scraper v2</span></div>
    </div>

    <nav class="nav">
      <div class="nav-item active" id="nav-scraper">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
        </svg>
        Scraper
      </div>
      <div class="nav-item" id="nav-leads">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
        </svg>
        Leads
      </div>
      <div class="nav-item" id="nav-whatsapp">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>
        </svg>
        WhatsApp
      </div>
    </nav>

    <div class="sidebar-footer">
      <p id="result-summary">No data yet</p>
    </div>
  </aside>

  <!-- Main -->
  <div class="main">
    <!-- Topbar -->
    <div class="topbar">
      <div class="topbar-form">
        <div class="input-wrap">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
          </svg>
          <input type="text" id="keyword" class="keyword-input" placeholder="e.g. coffee shop Bandung, toko aki Semarang..." autocomplete="off">
        </div>
        <div class="limit-wrap">
          <label for="limit">Limit</label>
          <input type="number" id="limit" class="limit-input" value="50" min="1" max="500">
        </div>
        <div class="headful-wrap">
          <label for="headful-toggle">Show browser</label>
          <label class="toggle">
            <input type="checkbox" id="headful-toggle">
            <span class="toggle-track"></span>
            <span class="toggle-thumb"></span>
          </label>
        </div>
        <button class="btn-scrape" id="btn-scrape" onclick="startScrape()">Start Scraping</button>
      </div>
    </div>

    <!-- Status bar -->
    <div class="statusbar">
      <div class="status-dot" id="status-dot"></div>
      <span id="status-text">Ready. Enter a keyword and start scraping.</span>
    </div>

    <!-- Content -->
    <div class="content">
      <div id="loading-bar" class="loading-bar" style="display:none"></div>

      <!-- Filter bar -->
      <div class="filter-bar">
        <button class="filter-chip active" data-tier="all" onclick="setFilter('all', this)">All</button>
        <button class="filter-chip chip-high" data-tier="HIGH POTENTIAL" onclick="setFilter('HIGH POTENTIAL', this)">High Potential</button>
        <button class="filter-chip chip-standard" data-tier="STANDARD" onclick="setFilter('STANDARD', this)">Standard</button>
        <button class="filter-chip chip-skip" data-tier="SKIP" onclick="setFilter('SKIP', this)">Skip</button>
        <div class="filter-right">
          <input type="text" class="search-input" id="search-box" placeholder="Filter by name or phone..." oninput="renderTable()">
          <span class="count-badge" id="count-badge">0 results</span>
        </div>
      </div>

      <!-- Table -->
      <div class="table-wrap">
        <table id="leads-table">
          <thead>
            <tr>
              <th onclick="sortBy('name')">Business Name <span class="sort-icon">↕</span></th>
              <th onclick="sortBy('rating')">Rating <span class="sort-icon">↕</span></th>
              <th onclick="sortBy('reviews')">Reviews <span class="sort-icon">↕</span></th>
              <th>Category</th>
              <th>Phone</th>
              <th>Website</th>
              <th onclick="sortBy('tier')">Tier <span class="sort-icon">↕</span></th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody id="table-body">
            <tr>
              <td colspan="8">
                <div class="empty-state">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
                  </svg>
                  <p>No results yet</p>
                  <small>Enter a keyword above and click Start Scraping</small>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>

<!-- WhatsApp Modal -->
<div class="modal-overlay" id="wa-modal">
  <div class="modal">
    <div class="modal-header">
      <span class="modal-title">Compose WhatsApp Message</span>
      <button class="modal-close" onclick="closeModal()">&times;</button>
    </div>
    <div>
      <label>Recipient: <span id="wa-name" style="color:var(--text);font-family:'Geist',sans-serif"></span></label>
    </div>
    <div>
      <label>Message</label>
      <textarea id="wa-message" placeholder="Type your message..."></textarea>
    </div>
    <div class="modal-footer">
      <button class="btn-ghost" onclick="closeModal()">Cancel</button>
      <button class="btn-send" onclick="sendWhatsApp()">Open WhatsApp</button>
    </div>
  </div>
</div>

<script>
  let allLeads = [];
  let activeFilter = 'all';
  let sortField = 'tier';
  let sortAsc = true;
  let pollTimer = null;
  let currentWaPhone = '';

  const WA_TEMPLATE = (name) =>
    "Halo " + name + ", perkenalkan saya dari tim kami.\n\n" +
    "Kami melihat bisnis Anda di Google Maps dan ingin menawarkan:\n" +
    "- Website profesional dengan landing page\n" +
    "- Sistem manajemen stok digital\n" +
    "- Invoice digital & laporan penjualan\n\n" +
    "Apakah Anda tertarik untuk mengetahui lebih lanjut? " +
    "Kami siap berdiskusi kapan saja.";

  // ── Scrape ────────────────────────────────────────────────────────────────
  async function startScrape() {
    const keyword = document.getElementById('keyword').value.trim();
    const limit   = parseInt(document.getElementById('limit').value) || 50;
    const headful = document.getElementById('headful-toggle').checked;

    if (!keyword) {
      document.getElementById('keyword').focus();
      return;
    }

    document.getElementById('btn-scrape').disabled = true;
    setStatus('running', 'Starting scraper for: ' + keyword);
    document.getElementById('loading-bar').style.display = 'block';
    allLeads = [];
    renderTable();

    try {
      const res = await fetch('/api/scrape', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ keyword, limit, headful, delay: 1.5 })
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Failed to start');
      pollStatus();
    } catch (e) {
      setStatus('error', 'Error: ' + e.message);
      document.getElementById('btn-scrape').disabled = false;
      document.getElementById('loading-bar').style.display = 'none';
    }
  }

  function pollStatus() {
    clearInterval(pollTimer);
    pollTimer = setInterval(async () => {
      try {
        const res = await fetch('/api/status');
        const data = await res.json();

        if (data.error) {
          setStatus('error', 'Error: ' + data.error);
          stopPoll();
          return;
        }

        if (data.running) {
          setStatus('running', 'Scraping "' + data.keyword + '"... (browser working)');
          return;
        }

        if (data.done) {
          await loadResults();
          setStatus('done', 'Done — ' + data.total + ' results for "' + data.keyword + '"');
          stopPoll();
        }
      } catch {}
    }, 2000);
  }

  function stopPoll() {
    clearInterval(pollTimer);
    document.getElementById('btn-scrape').disabled = false;
    document.getElementById('loading-bar').style.display = 'none';
  }

  async function loadResults() {
    const res = await fetch('/api/results');
    allLeads = await res.json();

    // Sort by tier order by default
    const order = {'HIGH POTENTIAL': 0, 'STANDARD': 1, 'SKIP': 2};
    allLeads.sort((a, b) => (order[a.tier] ?? 9) - (order[b.tier] ?? 9) || b.reviews - a.reviews);

    const high = allLeads.filter(l => l.tier === 'HIGH POTENTIAL').length;
    const std  = allLeads.filter(l => l.tier === 'STANDARD').length;
    const skip = allLeads.filter(l => l.tier === 'SKIP').length;
    document.getElementById('result-summary').textContent =
      'HIGH ' + high + ' / STD ' + std + ' / SKIP ' + skip;

    renderTable();
  }

  // ── Filter / Sort / Render ────────────────────────────────────────────────
  function setFilter(tier, el) {
    activeFilter = tier;
    document.querySelectorAll('.filter-chip').forEach(c => c.classList.remove('active'));
    el.classList.add('active');
    renderTable();
  }

  function sortBy(field) {
    if (sortField === field) { sortAsc = !sortAsc; }
    else { sortField = field; sortAsc = true; }
    document.querySelectorAll('th').forEach(th => th.classList.remove('sorted'));
    const thMap = { name:0, rating:1, reviews:2, tier:6 };
    const idx = thMap[field];
    if (idx !== undefined) document.querySelectorAll('th')[idx].classList.add('sorted');
    renderTable();
  }

  function renderTable() {
    const query = (document.getElementById('search-box').value || '').toLowerCase();
    let leads = allLeads;

    if (activeFilter !== 'all') leads = leads.filter(l => l.tier === activeFilter);
    if (query) leads = leads.filter(l =>
      l.name.toLowerCase().includes(query) ||
      l.phone.toLowerCase().includes(query) ||
      l.address.toLowerCase().includes(query)
    );

    // Sort
    const tierOrder = {'HIGH POTENTIAL': 0, 'STANDARD': 1, 'SKIP': 2};
    leads = [...leads].sort((a, b) => {
      let va, vb;
      if (sortField === 'name')    { va = a.name; vb = b.name; }
      else if (sortField === 'rating')  { va = a.rating; vb = b.rating; }
      else if (sortField === 'reviews') { va = a.reviews; vb = b.reviews; }
      else if (sortField === 'tier')    { va = tierOrder[a.tier] ?? 9; vb = tierOrder[b.tier] ?? 9; }
      else return 0;
      if (va < vb) return sortAsc ? -1 : 1;
      if (va > vb) return sortAsc ?  1 : -1;
      return 0;
    });

    document.getElementById('count-badge').textContent = leads.length + ' results';

    const tbody = document.getElementById('table-body');
    if (leads.length === 0) {
      tbody.innerHTML = '<tr><td colspan="8"><div class="empty-state">' +
        '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/></svg>' +
        '<p>No results match this filter</p><small>Try a different filter or search term</small></div></td></tr>';
      return;
    }

    tbody.innerHTML = leads.map(l => {
      const stars = renderStars(l.rating);
      const tierKey = l.tier.split(' ')[0];
      const tierLabel = tierKey === 'HIGH' ? 'HIGH' : tierKey;
      const waBtn = (l.tier === 'HIGH POTENTIAL' && l.phone)
        ? '<button class="btn-wa" onclick="openModal(\'' + escHtml(l.name) + '\',\'' + escHtml(l.phone) + '\')">' +
          '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347z"/><path d="M11.998 0C5.372 0 0 5.373 0 12c0 2.117.554 4.107 1.523 5.832L0 24l6.335-1.652A11.954 11.954 0 0011.998 24C18.626 24 24 18.627 24 12S18.626 0 11.998 0zm0 21.818a9.818 9.818 0 01-5.01-1.374l-.36-.213-3.76.984.998-3.65-.233-.375A9.82 9.82 0 012.18 12c0-5.42 4.4-9.818 9.818-9.818 5.42 0 9.82 4.398 9.82 9.818 0 5.42-4.4 9.818-9.82 9.818z"/></svg>' +
          'WhatsApp</button>'
        : '';

      return '<tr>' +
        '<td class="cell-name"><a href="' + escHtml(l.maps_url) + '" target="_blank">' + escHtml(l.name) + '</a>' +
        (l.address ? '<div class="sub">' + escHtml(l.address.substring(0, 60)) + (l.address.length > 60 ? '...' : '') + '</div>' : '') +
        '</td>' +
        '<td><span class="stars">' + stars + '</span><span class="rating-num">' + (l.rating > 0 ? l.rating.toFixed(1) : '—') + '</span></td>' +
        '<td class="reviews-num">' + (l.reviews > 0 ? l.reviews.toLocaleString() : '—') + '</td>' +
        '<td>' + escHtml(l.category || '—') + '</td>' +
        '<td style="font-family:\'Geist Mono\',monospace;font-size:12px">' + escHtml(l.phone || '—') + '</td>' +
        '<td>' + (l.has_website
          ? '<span class="chip chip-hassite">has site</span>'
          : '<span class="chip chip-nosite">no site</span>') + '</td>' +
        '<td><span class="tier-badge tier-' + tierKey + '" data-tip="' + escHtml(l.reason) + '">' + tierLabel + '</span></td>' +
        '<td>' + waBtn + '</td>' +
        '</tr>';
    }).join('');
  }

  function renderStars(rating) {
    if (!rating) return '—';
    const full  = Math.floor(rating);
    const half  = rating - full >= 0.25 && rating - full < 0.75 ? 1 : 0;
    const empty = 5 - full - half;
    return '★'.repeat(full) + (half ? '½' : '') + '☆'.repeat(Math.max(0, empty));
  }

  // ── Status ────────────────────────────────────────────────────────────────
  function setStatus(type, msg) {
    const dot = document.getElementById('status-dot');
    dot.className = 'status-dot ' + type;
    document.getElementById('status-text').textContent = msg;
  }

  // ── WhatsApp Modal ─────────────────────────────────────────────────────────
  function openModal(name, phone) {
    currentWaPhone = phone;
    document.getElementById('wa-name').textContent = name + ' — ' + phone;
    document.getElementById('wa-message').value = WA_TEMPLATE(name);
    document.getElementById('wa-modal').classList.add('open');
  }

  function closeModal() {
    document.getElementById('wa-modal').classList.remove('open');
  }

  function sendWhatsApp() {
    const msg = document.getElementById('wa-message').value.trim();
    const phone = currentWaPhone.replace(/[^0-9]/g, '');
    const intl = phone.startsWith('0') ? '62' + phone.slice(1) : phone;
    const url = 'https://wa.me/' + intl + '?text=' + encodeURIComponent(msg);
    window.open(url, '_blank');
    closeModal();
  }

  // Close modal on overlay click
  document.getElementById('wa-modal').addEventListener('click', function(e) {
    if (e.target === this) closeModal();
  });

  // Enter to scrape
  document.getElementById('keyword').addEventListener('keydown', function(e) {
    if (e.key === 'Enter') startScrape();
  });

  // Escape to close modal
  document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') closeModal();
  });

  function escHtml(s) {
    if (!s) return '';
    return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
  }
</script>
</body>
</html>`
