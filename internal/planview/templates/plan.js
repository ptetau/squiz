// squiz-plan client-side controller.
//
// Wires the tabbed plan view: tab switching, cross-ref badge navigation,
// per-item feedback (status / note / inline edits), per-item options
// chooser (v0.4.0), per-section notes (v0.4.0), plan-level note
// (v0.4.0), proposed-item form (v0.4.0), and the copy-json modal.
// Expects two globals injected by the template at parse time:
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

    // In-memory feedback state.
    //   feedback[itemId]    = { status, note, edits, chose }
    //   sectionNotes[id]    = string
    //   planNote            = string
    //   proposedItems[]     = { section, title, desc, refs[] }
    const state = {
      feedback: {},
      sectionNotes: {},
      planNote: '',
      proposedItems: []
    };

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
      const st = e.state;
      if (st && st.kind === 'item') {
        navigateToItem(st.itemId, { smooth: false });
      } else if (st && st.kind === 'tab') {
        activateTab(st.tabId);
      } else {
        // No state — interpret current hash.
        applyHash(location.hash, { smooth: false });
      }
    });

    // ── FEEDBACK WIDGETS ────────────────────────────────────────────
    function ensureFb(itemId) {
      if (!state.feedback[itemId]) {
        state.feedback[itemId] = { status: null, note: '', edits: {}, chose: null };
      }
      return state.feedback[itemId];
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

    // ── PER-ITEM OPTIONS CHOOSER (v0.4.0) ───────────────────────────
    function selectOption(btn) {
      const group = btn.closest('.item-options');
      if (!group) return;
      const itemId = group.dataset.item;
      const optionId = btn.dataset.option;
      ensureFb(itemId).chose = optionId;
      group.querySelectorAll('.item-option').forEach(o => {
        const isMe = o === btn;
        o.classList.toggle('selected', isMe);
        o.setAttribute('aria-checked', isMe ? 'true' : 'false');
        o.setAttribute('tabindex', isMe ? '0' : '-1');
      });
      updateStats();
    }

    document.querySelectorAll('.item-options').forEach(group => {
      const opts = Array.from(group.querySelectorAll('.item-option'));
      opts.forEach(btn => {
        btn.addEventListener('click', () => selectOption(btn));
      });
      group.addEventListener('keydown', e => {
        const current = document.activeElement;
        const idx = opts.indexOf(current);
        if (idx < 0) return;
        let next = null;
        switch (e.key) {
          case 'ArrowRight':
          case 'ArrowDown':
            next = opts[(idx + 1) % opts.length]; break;
          case 'ArrowLeft':
          case 'ArrowUp':
            next = opts[(idx - 1 + opts.length) % opts.length]; break;
          case 'Home':
            next = opts[0]; break;
          case 'End':
            next = opts[opts.length - 1]; break;
          case ' ':
          case 'Enter':
            e.preventDefault();
            selectOption(current);
            return;
          default:
            return;
        }
        if (next) {
          e.preventDefault();
          opts.forEach(o => o.setAttribute('tabindex', o === next ? '0' : '-1'));
          next.focus();
        }
      });
    });

    // ── SECTION NOTES (v0.4.0) ──────────────────────────────────────
    document.querySelectorAll('.section-note').forEach(ta => {
      ta.addEventListener('input', () => {
        const sid = ta.dataset.section;
        state.sectionNotes[sid] = ta.value;
        ta.style.height = 'auto';
        ta.style.height = ta.scrollHeight + 'px';
        updateStats();
      });
    });

    // ── PLAN-LEVEL NOTE (v0.4.0) ────────────────────────────────────
    const planNoteEl = document.querySelector('.plan-note');
    if (planNoteEl) {
      planNoteEl.addEventListener('input', () => {
        state.planNote = planNoteEl.value;
        // Live-refresh the JSON preview while the modal is open.
        refreshModalPreview();
        updateStats();
      });
    }

    // ── ADD-ITEM FORM (v0.4.0) ──────────────────────────────────────
    // Populate the refs <select> on every form with all known item IDs.
    function populateRefsSelect(sel) {
      sel.innerHTML = '';
      Object.keys(itemIndex).forEach(id => {
        const opt = document.createElement('option');
        opt.value = id;
        opt.textContent = id;
        sel.appendChild(opt);
      });
    }

    document.querySelectorAll('.add-item-form').forEach(form => {
      const sel = form.querySelector('.add-item-refs');
      if (sel) populateRefsSelect(sel);
    });

    document.querySelectorAll('.add-item-btn').forEach(btn => {
      const sectionId = btn.dataset.section;
      const form = document.getElementById('add-item-form-' + sectionId);
      if (!form) return;
      btn.addEventListener('click', () => {
        const open = !form.hidden;
        form.hidden = open;
        btn.setAttribute('aria-expanded', open ? 'false' : 'true');
        if (!open) {
          const titleInput = form.querySelector('.add-item-title');
          if (titleInput) {
            try { titleInput.focus(); } catch (e) { /* noop */ }
          }
        }
      });
      const cancelBtn = form.querySelector('.add-item-cancel');
      if (cancelBtn) {
        cancelBtn.addEventListener('click', () => {
          form.reset();
          form.hidden = true;
          btn.setAttribute('aria-expanded', 'false');
          try { btn.focus(); } catch (e) { /* noop */ }
        });
      }
      form.addEventListener('submit', e => {
        e.preventDefault();
        const titleInput = form.querySelector('.add-item-title');
        const descInput = form.querySelector('.add-item-desc');
        const refsSel = form.querySelector('.add-item-refs');
        const title = (titleInput && titleInput.value || '').trim();
        const desc = (descInput && descInput.value || '').trim();
        const refs = refsSel
          ? Array.from(refsSel.selectedOptions).map(o => o.value)
          : [];
        if (!title && !desc) return;
        const proposal = { section: sectionId, title: title, desc: desc, refs: refs };
        state.proposedItems.push(proposal);
        renderProposed(sectionId, proposal, state.proposedItems.length - 1);
        form.reset();
        form.hidden = true;
        btn.setAttribute('aria-expanded', 'false');
        try { btn.focus(); } catch (e) { /* noop */ }
        updateStats();
      });
    });

    function renderProposed(sectionId, proposal, idx) {
      const wrap = document.querySelector('.proposed-items[data-section="' + sectionId + '"]');
      if (!wrap) return;
      const art = document.createElement('article');
      art.className = 'item proposed';
      art.dataset.proposedIdx = String(idx);
      art.setAttribute('aria-label', 'Proposed item: ' + (proposal.title || '(untitled)'));

      const header = document.createElement('header');
      header.className = 'item-head';
      const h2 = document.createElement('h2');
      h2.className = 'item-title';
      h2.textContent = proposal.title || '(untitled)';
      const badge = document.createElement('span');
      badge.className = 'item-id proposed-badge';
      badge.textContent = 'PROPOSED';
      header.appendChild(h2);
      header.appendChild(badge);
      art.appendChild(header);

      if (proposal.desc) {
        const p = document.createElement('p');
        p.className = 'item-desc';
        p.textContent = proposal.desc;
        art.appendChild(p);
      }

      if (proposal.refs && proposal.refs.length) {
        const ul = document.createElement('ul');
        ul.className = 'refs';
        ul.setAttribute('aria-label', 'References');
        const lbl = document.createElement('li');
        lbl.className = 'refs-label';
        lbl.textContent = 'refs';
        ul.appendChild(lbl);
        proposal.refs.forEach(rid => {
          const li = document.createElement('li');
          const a = document.createElement('a');
          a.href = '#item-' + rid;
          a.dataset.target = rid;
          a.textContent = rid;
          a.addEventListener('click', ev => {
            ev.preventDefault();
            if (itemIndex[rid]) navigateToItem(rid);
          });
          li.appendChild(a);
          ul.appendChild(li);
        });
        art.appendChild(ul);
      }

      const actions = document.createElement('div');
      actions.className = 'proposed-actions';
      const rm = document.createElement('button');
      rm.type = 'button';
      rm.className = 'proposed-remove';
      rm.textContent = 'remove';
      rm.addEventListener('click', () => {
        // Null out (preserves indexes used as dataset.proposedIdx) and
        // remove the DOM node; buildExport filters nulls.
        const i = Number(art.dataset.proposedIdx);
        if (!Number.isNaN(i)) state.proposedItems[i] = null;
        art.remove();
        updateStats();
      });
      actions.appendChild(rm);
      art.appendChild(actions);

      wrap.appendChild(art);
    }

    // ── EXPORT BUILDER ──────────────────────────────────────────────
    function buildExport() {
      // Per-item feedback: drop empties, diff edits against original.
      const entries = [];
      let approved = 0, questioned = 0, rejected = 0;
      let withNotes = 0, withEdits = 0, withChose = 0;
      Object.keys(state.feedback).forEach(id => {
        const f = state.feedback[id];
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
        const chose = f.chose || null;
        if (!f.status && !note && !hasEdits && !chose) return;
        entries.push({
          id: id,
          status: f.status || null,
          anchor: '#item-' + id,
          note: note || null,
          edits: hasEdits ? edits : null,
          chose: chose
        });
        if (f.status === 'approved') approved++;
        else if (f.status === 'questioned') questioned++;
        else if (f.status === 'rejected') rejected++;
        if (note) withNotes++;
        if (hasEdits) withEdits++;
        if (chose) withChose++;
      });

      // Section notes: drop empties.
      const sectionNotes = {};
      let sectionsWithNotes = 0;
      Object.keys(state.sectionNotes).forEach(sid => {
        const v = (state.sectionNotes[sid] || '').trim();
        if (v) {
          sectionNotes[sid] = v;
          sectionsWithNotes++;
        }
      });

      const planNote = (state.planNote || '').trim();

      // Proposed items: skip nulls (removed).
      const proposed = state.proposedItems
        .filter(p => p && (p.title || p.desc || (p.refs && p.refs.length)));

      const total = Object.keys(itemIndex).length;
      const out = {
        plan: PLAN.title || '',
        source: SOURCE,
        generatedAt: new Date().toISOString().slice(0, 19) + 'Z',
        feedback: entries
      };
      if (Object.keys(sectionNotes).length) {
        out.section_notes = sectionNotes;
      }
      if (planNote) {
        out.plan_note = planNote;
      }
      out.proposed_items = proposed;
      out.summary = {
        total: total,
        approved: approved,
        questioned: questioned,
        rejected: rejected,
        withNotes: withNotes,
        withEdits: withEdits,
        withChose: withChose,
        sectionsWithNotes: sectionsWithNotes,
        hasPlanNote: !!planNote,
        proposedItems: proposed.length
      };
      return out;
    }

    function updateStats() {
      const totalEl = document.querySelector('.copy-bar-total');
      const touchedEl = document.querySelector('.copy-bar-touched');
      const btn = document.querySelector('.copy-bar-btn');
      const total = Object.keys(itemIndex).length;
      let touched = 0;
      Object.keys(state.feedback).forEach(id => {
        const f = state.feedback[id];
        const note = (f.note || '').trim();
        const hasEdits = f.edits && (
          (typeof f.edits.title === 'string' && f.edits.title.length > 0) ||
          (typeof f.edits.desc === 'string' && f.edits.desc.length > 0)
        );
        if (f.status || note || hasEdits || f.chose) touched++;
      });
      // The copy-bar counter tracks per-item touched only (its label
      // reads "items touched"); section/plan/proposed inputs separately
      // enable the export button below.
      const sectionsWithNotes = Object.keys(state.sectionNotes)
        .filter(k => (state.sectionNotes[k] || '').trim()).length;
      const proposedCount = state.proposedItems.filter(p => p).length;
      const planNoteLen = (state.planNote || '').trim().length;

      if (totalEl) totalEl.textContent = String(total).padStart(2, '0');
      if (touchedEl) touchedEl.textContent = String(touched).padStart(2, '0');
      if (btn) {
        btn.disabled = (touched + sectionsWithNotes + proposedCount + (planNoteLen ? 1 : 0)) === 0;
      }
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

    function refreshModalPreview() {
      if (!modal || !modalPre) return;
      if (modal.style.display === 'none' || modal.style.display === '') return;
      const json = JSON.stringify(buildExport(), null, 2);
      modalPre.innerHTML = syntaxHighlight(json);
      modalPre.dataset.raw = json;
    }

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
        // Focus trap: cycle between close + plan-note + copy.
        const focusables = [closeBtn, planNoteEl, copyBtn].filter(Boolean);
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
        // Rebuild on copy so a late edit in plan-note is captured.
        const json = JSON.stringify(buildExport(), null, 2);
        if (modalPre) {
          modalPre.innerHTML = syntaxHighlight(json);
          modalPre.dataset.raw = json;
        }
        const raw = json;
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
