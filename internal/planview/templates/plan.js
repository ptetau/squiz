// squiz-plan client-side controller.
//
// Wires the tabbed plan view: tab switching, cross-ref badge navigation,
// per-item feedback (status / note / inline edits), and the copy-json
// modal. Expects two globals injected by the template at parse time:
//
//   window.PLAN   — { title, lede, sections: [{id, label, items: [...]}], ... }
//   window.SOURCE — { file, basename }
//
// No build step. No deps. ES2017+ (arrow funcs, template literals,
// async/await). Bind on DOMContentLoaded so the template can ship the
// <script> anywhere in <body>.
(function () {
  'use strict';

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', boot);
  } else {
    boot();
  }

  function boot() {
    const PLAN = window.PLAN || { sections: [] };
    const SOURCE = window.SOURCE || { file: '', basename: '' };

    // Build quick lookup tables.
    // itemIndex: itemId -> { sectionId, item }
    const itemIndex = {};
    (PLAN.sections || []).forEach(s => {
      (s.items || []).forEach(it => {
        itemIndex[it.id] = { sectionId: s.id, item: it };
      });
    });

    // In-memory feedback state: feedback[itemId] = { status, note, edits }
    const feedback = {};

    // ── TAB SWITCHING ───────────────────────────────────────────────
    const tabs = Array.from(document.querySelectorAll('.plan-tab'));
    const panels = Array.from(document.querySelectorAll('.tabpanel'));

    function panelFor(tabId) {
      return document.getElementById('panel-' + tabId);
    }
    function tabFor(tabId) {
      return document.getElementById('tab-' + tabId);
    }

    function activateTab(tabId, opts) {
      opts = opts || {};
      const tab = tabFor(tabId);
      const panel = panelFor(tabId);
      if (!tab || !panel) return false;
      tabs.forEach(t => {
        const me = t === tab;
        t.setAttribute('aria-selected', me ? 'true' : 'false');
        t.setAttribute('tabindex', me ? '0' : '-1');
      });
      panels.forEach(p => {
        p.hidden = (p !== panel);
      });
      if (opts.focus) {
        tab.focus({ preventScroll: !!opts.preventScroll });
      }
      return true;
    }

    tabs.forEach(tab => {
      tab.addEventListener('click', () => {
        const tabId = tab.dataset.tab;
        if (!activateTab(tabId)) return;
        const newHash = '#tab-' + tabId;
        if (location.hash !== newHash) {
          history.pushState({ kind: 'tab', tabId: tabId }, '', newHash);
        }
      });
    });

    // Keyboard nav across the tablist.
    const tablist = document.querySelector('.plan-tabs');
    if (tablist) {
      tablist.addEventListener('keydown', e => {
        const idx = tabs.indexOf(document.activeElement);
        if (idx < 0) return;
        let next = null;
        switch (e.key) {
          case 'ArrowRight': next = tabs[(idx + 1) % tabs.length]; break;
          case 'ArrowLeft':  next = tabs[(idx - 1 + tabs.length) % tabs.length]; break;
          case 'Home':       next = tabs[0]; break;
          case 'End':        next = tabs[tabs.length - 1]; break;
          case 'Enter':
          case ' ':
            e.preventDefault();
            document.activeElement.click();
            return;
          default: return;
        }
        if (next) {
          e.preventDefault();
          activateTab(next.dataset.tab, { focus: true, preventScroll: true });
          const newHash = '#tab-' + next.dataset.tab;
          if (location.hash !== newHash) {
            history.pushState({ kind: 'tab', tabId: next.dataset.tab }, '', newHash);
          }
        }
      });
    }

    // ── CROSS-REF BADGE NAVIGATION ──────────────────────────────────
    let highlightTimer = null;
    function navigateToItem(itemId, opts) {
      opts = opts || {};
      const entry = itemIndex[itemId];
      if (!entry) return false;
      activateTab(entry.sectionId);
      const target = document.getElementById('item-' + itemId);
      if (!target) return false;
      // Scroll first; the sticky tab strip is handled via scroll-margin.
      target.scrollIntoView({
        behavior: opts.smooth === false ? 'auto' : 'smooth',
        block: 'start'
      });
      // Highlight + focus.
      target.classList.remove('is-target');
      // Force a reflow so the animation restarts when revisiting same id.
      void target.offsetWidth;
      target.classList.add('is-target');
      if (highlightTimer) clearTimeout(highlightTimer);
      highlightTimer = setTimeout(() => {
        target.classList.remove('is-target');
      }, 1500);
      setTimeout(() => {
        try { target.focus({ preventScroll: true }); } catch (e) { /* noop */ }
      }, 250);
      return true;
    }

    document.querySelectorAll('.refs a[data-target]').forEach(a => {
      a.addEventListener('click', e => {
        e.preventDefault();
        const itemId = a.dataset.target;
        if (!itemIndex[itemId]) return;
        navigateToItem(itemId);
        const newHash = '#item-' + itemId;
        if (location.hash !== newHash) {
          history.pushState({ kind: 'item', itemId: itemId }, '', newHash);
        }
      });
    });

    // popstate: respond to back / forward.
    window.addEventListener('popstate', e => {
      const state = e.state;
      if (state && state.kind === 'item') {
        navigateToItem(state.itemId, { smooth: false });
      } else if (state && state.kind === 'tab') {
        activateTab(state.tabId);
      } else {
        // No state — interpret current hash.
        applyHash(location.hash, { smooth: false });
      }
    });

    // ── FEEDBACK WIDGETS ────────────────────────────────────────────
    function ensureFb(itemId) {
      if (!feedback[itemId]) feedback[itemId] = { status: null, note: '', edits: {} };
      return feedback[itemId];
    }

    document.querySelectorAll('.feedback').forEach(fb => {
      const itemId = fb.dataset.item;
      fb.querySelectorAll('.fb-btn').forEach(btn => {
        btn.addEventListener('click', () => {
          const value = btn.dataset.value;
          const current = ensureFb(itemId);
          // Toggle off if the same status is re-clicked.
          if (current.status === value) {
            current.status = null;
            fb.removeAttribute('data-status');
            btn.setAttribute('aria-pressed', 'false');
          } else {
            current.status = value;
            fb.setAttribute('data-status', value);
            fb.querySelectorAll('.fb-btn').forEach(b => {
              b.setAttribute('aria-pressed', b === btn ? 'true' : 'false');
            });
          }
          updateStats();
        });
      });
      const note = fb.querySelector('.fb-note');
      if (note) {
        note.addEventListener('input', () => {
          ensureFb(itemId).note = note.value;
          note.style.height = 'auto';
          note.style.height = note.scrollHeight + 'px';
          updateStats();
        });
      }
      const editBtn = fb.querySelector('.fb-edit');
      if (editBtn) {
        editBtn.addEventListener('click', () => {
          const article = document.getElementById('item-' + itemId);
          if (!article) return;
          const titleEl = article.querySelector('.item-title');
          const descEl = article.querySelector('.item-desc');
          const on = editBtn.getAttribute('aria-pressed') !== 'true';
          editBtn.setAttribute('aria-pressed', on ? 'true' : 'false');
          editBtn.textContent = on ? 'done' : 'edit';
          [titleEl, descEl].forEach(el => {
            if (!el) return;
            if (on) {
              el.setAttribute('contenteditable', 'true');
              if (!el.dataset.bound) {
                el.dataset.bound = '1';
                el.addEventListener('input', () => {
                  const f = ensureFb(itemId);
                  const field = el.classList.contains('item-title') ? 'title' : 'desc';
                  f.edits[field] = el.textContent;
                  updateStats();
                });
              }
            } else {
              el.removeAttribute('contenteditable');
            }
          });
          if (on && titleEl) {
            try { titleEl.focus(); } catch (e) { /* noop */ }
          }
        });
      }
    });

    // ── EXPORT BUILDER ──────────────────────────────────────────────
    function buildExport() {
      // Filter empties + diff edits against original PLAN values.
      const entries = [];
      let withStatus = 0, withNotes = 0, withEdits = 0;
      Object.keys(feedback).forEach(id => {
        const f = feedback[id];
        const orig = (itemIndex[id] && itemIndex[id].item) || {};
        const edits = {};
        if (f.edits) {
          if (typeof f.edits.title === 'string' && f.edits.title.trim() !== '' && f.edits.title !== orig.title) {
            edits.title = f.edits.title;
          }
          if (typeof f.edits.desc === 'string' && f.edits.desc.trim() !== '' && f.edits.desc !== orig.desc) {
            edits.desc = f.edits.desc;
          }
        }
        const note = (f.note || '').trim();
        const hasEdits = Object.keys(edits).length > 0;
        if (!f.status && !note && !hasEdits) return;
        entries.push({
          id: id,
          status: f.status || null,
          note: note || null,
          edits: hasEdits ? edits : null
        });
        if (f.status) withStatus++;
        if (note) withNotes++;
        if (hasEdits) withEdits++;
      });
      const total = Object.keys(itemIndex).length;
      return {
        spec: PLAN.title || '',
        source: SOURCE,
        generatedAt: new Date().toISOString().slice(0, 19) + 'Z',
        feedback: entries,
        summary: {
          total: total,
          withStatus: withStatus,
          withNotes: withNotes,
          withEdits: withEdits
        }
      };
    }

    function updateStats() {
      const totalEl = document.querySelector('.copy-bar-total');
      const touchedEl = document.querySelector('.copy-bar-touched');
      const btn = document.querySelector('.copy-bar-btn');
      const total = Object.keys(itemIndex).length;
      let touched = 0;
      Object.keys(feedback).forEach(id => {
        const f = feedback[id];
        const note = (f.note || '').trim();
        const hasEdits = f.edits && (
          (typeof f.edits.title === 'string' && f.edits.title.length > 0) ||
          (typeof f.edits.desc === 'string' && f.edits.desc.length > 0)
        );
        if (f.status || note || hasEdits) touched++;
      });
      if (totalEl) totalEl.textContent = String(total).padStart(2, '0');
      if (touchedEl) touchedEl.textContent = String(touched).padStart(2, '0');
      if (btn) btn.disabled = touched === 0;
    }

    function syntaxHighlight(json) {
      return json
        .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
        .replace(/("(\\u[a-fA-F0-9]{4}|\\[^u]|[^\\"])*")(\s*:)?/g, function (m, p1, p2, p3) {
          if (p3) return '<span class="k">' + p1 + '</span>' + p3;
          return '<span class="s">' + p1 + '</span>';
        })
        .replace(/\b(true|false|null)\b/g, '<span class="b">$1</span>')
        .replace(/(?<!["\w-])(-?\d+(\.\d+)?)/g, '<span class="n">$1</span>');
    }

    // ── COPY-JSON MODAL ─────────────────────────────────────────────
    const modal = document.querySelector('.modal-overlay');
    const modalPre = modal ? modal.querySelector('pre') : null;
    const copyBarBtn = document.querySelector('.copy-bar-btn');
    const closeBtn = modal ? modal.querySelector('.x') : null;
    const copyBtn = modal ? modal.querySelector('.copy-btn') : null;
    let modalReturnFocus = null;

    function openModal() {
      if (!modal || !modalPre) return;
      if (copyBarBtn && copyBarBtn.disabled) return;
      const json = JSON.stringify(buildExport(), null, 2);
      modalPre.innerHTML = syntaxHighlight(json);
      modalPre.dataset.raw = json;
      modalReturnFocus = document.activeElement;
      modal.style.display = 'flex';
      modal.setAttribute('aria-hidden', 'false');
      setTimeout(() => { if (copyBtn) copyBtn.focus(); }, 0);
    }
    function closeModal() {
      if (!modal) return;
      modal.style.display = 'none';
      modal.setAttribute('aria-hidden', 'true');
      if (modalReturnFocus && typeof modalReturnFocus.focus === 'function') {
        modalReturnFocus.focus();
      }
    }

    if (copyBarBtn) copyBarBtn.addEventListener('click', openModal);
    if (closeBtn) closeBtn.addEventListener('click', closeModal);
    if (modal) {
      modal.addEventListener('click', e => { if (e.target === modal) closeModal(); });
    }
    document.addEventListener('keydown', e => {
      if (!modal || modal.style.display === 'none') return;
      if (e.key === 'Escape') { e.preventDefault(); closeModal(); return; }
      if (e.key === 'Tab') {
        const focusables = [closeBtn, copyBtn].filter(Boolean);
        if (focusables.length === 0) return;
        const idx = focusables.indexOf(document.activeElement);
        if (idx < 0) {
          e.preventDefault();
          focusables[0].focus();
          return;
        }
        e.preventDefault();
        const dir = e.shiftKey ? -1 : 1;
        focusables[(idx + dir + focusables.length) % focusables.length].focus();
      }
    });

    if (copyBtn) {
      copyBtn.addEventListener('click', async () => {
        const raw = (modalPre && modalPre.dataset.raw) || '';
        try {
          await navigator.clipboard.writeText(raw);
        } catch (err) {
          const ta = document.createElement('textarea');
          ta.value = raw;
          document.body.appendChild(ta);
          ta.select();
          try { document.execCommand('copy'); } catch (e2) { /* noop */ }
          document.body.removeChild(ta);
        }
        copyBtn.classList.add('copied');
        const orig = copyBtn.textContent;
        copyBtn.textContent = 'copied';
        setTimeout(() => {
          copyBtn.classList.remove('copied');
          copyBtn.textContent = orig;
        }, 2000);
      });
    }

    // ── INITIAL HASH ROUTING ────────────────────────────────────────
    function applyHash(hash, opts) {
      opts = opts || {};
      if (!hash || hash === '#') {
        // Default: first tab.
        if (tabs.length) activateTab(tabs[0].dataset.tab);
        return;
      }
      const h = hash.replace(/^#/, '');
      if (h.indexOf('item-') === 0) {
        const id = h.slice('item-'.length);
        if (itemIndex[id]) {
          navigateToItem(id, { smooth: opts.smooth });
          return;
        }
      }
      if (h.indexOf('tab-') === 0) {
        const tabId = h.slice('tab-'.length);
        if (activateTab(tabId)) return;
      }
      // Fall back to first tab.
      if (tabs.length) activateTab(tabs[0].dataset.tab);
    }

    applyHash(location.hash, { smooth: false });
    updateStats();
  }
})();
