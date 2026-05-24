---
name: squiz-plan
description: >
  Sibling to /squiz. Takes a structured plan or spec (overview, functional
  requirements, non-functional requirements, cases, engineering requirements,
  build steps) and renders it as a tabbed, cross-referenced interactive HTML
  document with the Apple //e × IBM Plex aesthetic. Every layer links back
  to the one above so "why does this step exist?" is one click away.
  STATUS: skeleton — full implementation lands in v0.3.0.
---

# squiz-plan: structured plans, visible threads

**Status: placeholder skeleton.** The `squiz-plan` binary is not yet implemented; this file ships so the install scripts can lay down the skills directory ahead of the v0.3.0 cut. Once `squiz-plan` lands, this document will describe the JSON schema, tab layout, cross-reference ID convention (`OVR-…`, `FR-…`, `NFR-…`, `CASE-…`, `ENG-…`, `BUILD-…`), and the full feedback export shape.

For now, use [[squiz]] for clarification flows.

## Planned shape (so you can pre-author plans against it)

```
plan/
├── index.json              # section order + which files to include
├── overview.json           # items with id: OVR-…
├── functional.json         # FR-…
├── non-functional.json     # NFR-…
├── cases.json              # CASE-…
├── engineering.json        # ENG-…
└── build.json              # BUILD-… (often per-component, nested)
```

Each item carries `id`, `title`, `desc`, optional `art` (same forms as squiz: `wf:`, DSL, raw SVG, `none`), and `refs: ["OVR-1", "FR-3"]` — parent references rendered as clickable badges that switch tabs and highlight the target.

Run it (when shipped):

```bash
squiz-plan plan/index.json
```
