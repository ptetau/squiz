package renderer

// ArchLibrary — system-design icon registry. Agents reference these via
// `"art": "arch:<name>"`. Each entry is a `<g>...</g>` SVG fragment sized
// to roughly fit a 40×40 box centered on (50,30) within the standard
// viewBox 0 0 100 60. The dispatcher (see art.go) wraps each `<g>` in
// `<svg viewBox='0 0 100 60' style='width:55%;height:auto'>` so they
// render at consistent visual size alongside `wf:*` entries.
//
// Style conventions (mirroring wf.go):
//   - minimal line-art, single-weight strokes (1.4–1.6)
//   - primary strokes in var(--accent), secondary detail in var(--ink)
//   - fills are 'none' or var(--accent-soft)/var(--accent) where useful
//   - feel "retro / clean" — they sit beside the wireframes
//
// Categories (30 entries total):
//   Compute (6):       server, container, pod, function, worker, scheduler
//   Data (6):          database, table, blob, storage, cache, stream
//   Network (7):       load-balancer, gateway, cdn, dns, vpc, subnet, firewall
//   Services (3):      api, queue, topic
//   Observability (3): log, metric, trace
//   Identity (3):      user, mobile, browser
//   Security (2):      secret, key-icon
var archLibrary = map[string]string{

	// ── Compute ─────────────────────────────────────────────────────
	// server: stacked rackmount unit
	"server": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='30' y='14' width='40' height='10' rx='1'/>` +
		`<rect x='30' y='26' width='40' height='10' rx='1'/>` +
		`<rect x='30' y='38' width='40' height='10' rx='1'/>` +
		`<circle cx='36' cy='19' r='1.2' fill='var(--accent)' stroke='none'/>` +
		`<circle cx='36' cy='31' r='1.2' fill='var(--accent)' stroke='none'/>` +
		`<circle cx='36' cy='43' r='1.2' fill='var(--accent)' stroke='none'/>` +
		`<line x1='42' y1='19' x2='64' y2='19' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='42' y1='31' x2='64' y2='31' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='42' y1='43' x2='64' y2='43' stroke='var(--ink)' stroke-width='1'/>` +
		`</g>`,

	// container: a "shipping container" cube with ridges
	"container": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='30' y='16' width='40' height='28'/>` +
		`<line x1='36' y1='16' x2='36' y2='44' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='44' y1='16' x2='44' y2='44' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='52' y1='16' x2='52' y2='44' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='60' y1='16' x2='60' y2='44' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='30' y1='22' x2='70' y2='22'/>` +
		`</g>`,

	// pod: hexagon (k8s vibe)
	"pod": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<polygon points='50,12 70,22 70,38 50,48 30,38 30,22' fill='var(--accent-soft)'/>` +
		`<line x1='50' y1='12' x2='50' y2='30' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='30' y1='22' x2='50' y2='30' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='70' y1='22' x2='50' y2='30' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='30' y1='38' x2='50' y2='30' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='70' y1='38' x2='50' y2='30' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='50' y1='48' x2='50' y2='30' stroke='var(--ink)' stroke-width='1'/>` +
		`</g>`,

	// function: lambda glyph in rounded square
	"function": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='32' y='14' width='36' height='32' rx='4'/>` +
		`<path d='M 42 40 L 50 24 L 46 18 M 50 24 L 58 40' stroke='var(--accent)' stroke-width='1.6' stroke-linecap='round'/>` +
		`</g>`,

	// worker: gear
	"worker": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<circle cx='50' cy='30' r='10'/>` +
		`<circle cx='50' cy='30' r='4' fill='var(--accent)' stroke='none'/>` +
		`<g stroke='var(--accent)' stroke-width='1.5' stroke-linecap='round'>` +
		`<line x1='50' y1='14' x2='50' y2='18'/>` +
		`<line x1='50' y1='42' x2='50' y2='46'/>` +
		`<line x1='34' y1='30' x2='38' y2='30'/>` +
		`<line x1='62' y1='30' x2='66' y2='30'/>` +
		`<line x1='38' y1='18' x2='40.8' y2='20.8'/>` +
		`<line x1='59.2' y1='39.2' x2='62' y2='42'/>` +
		`<line x1='62' y1='18' x2='59.2' y2='20.8'/>` +
		`<line x1='40.8' y1='39.2' x2='38' y2='42'/>` +
		`</g></g>`,

	// scheduler: clock + bracket (scheduled jobs)
	"scheduler": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<circle cx='50' cy='30' r='14'/>` +
		`<line x1='50' y1='30' x2='50' y2='20' stroke='var(--accent)' stroke-width='1.6' stroke-linecap='round'/>` +
		`<line x1='50' y1='30' x2='58' y2='34' stroke='var(--accent)' stroke-width='1.6' stroke-linecap='round'/>` +
		`<circle cx='50' cy='30' r='1.4' fill='var(--accent)' stroke='none'/>` +
		`<path d='M 32 18 L 28 18 L 28 42 L 32 42' stroke='var(--ink)' stroke-width='1'/>` +
		`<path d='M 68 18 L 72 18 L 72 42 L 68 42' stroke='var(--ink)' stroke-width='1'/>` +
		`</g>`,

	// ── Data ────────────────────────────────────────────────────────
	// database: classic cylinder
	"database": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<ellipse cx='50' cy='16' rx='16' ry='4' fill='var(--accent-soft)'/>` +
		`<path d='M 34 16 L 34 44 A 16 4 0 0 0 66 44 L 66 16'/>` +
		`<path d='M 34 26 A 16 4 0 0 0 66 26' stroke='var(--ink)' stroke-width='1'/>` +
		`<path d='M 34 36 A 16 4 0 0 0 66 36' stroke='var(--ink)' stroke-width='1'/>` +
		`</g>`,

	// table: ruled rows + header
	"table": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='28' y='14' width='44' height='32'/>` +
		`<rect x='28' y='14' width='44' height='6' fill='var(--accent-soft)' stroke='var(--accent)' stroke-width='1.4'/>` +
		`<line x1='28' y1='28' x2='72' y2='28' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='28' y1='36' x2='72' y2='36' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='44' y1='20' x2='44' y2='46' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='58' y1='20' x2='58' y2='46' stroke='var(--ink)' stroke-width='1'/>` +
		`</g>`,

	// blob: amorphous data lump
	"blob": `<g fill='var(--accent-soft)' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<path d='M 36 18 Q 28 24 32 34 Q 30 44 42 46 Q 54 50 62 42 Q 72 38 68 26 Q 66 16 54 16 Q 44 14 36 18 Z'/>` +
		`<circle cx='44' cy='28' r='1.6' fill='var(--ink)' stroke='none'/>` +
		`<circle cx='56' cy='32' r='1.6' fill='var(--ink)' stroke='none'/>` +
		`<circle cx='48' cy='38' r='1.6' fill='var(--ink)' stroke='none'/>` +
		`</g>`,

	// storage: disk/drive stack
	"storage": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='30' y='16' width='40' height='12' rx='1'/>` +
		`<rect x='30' y='32' width='40' height='12' rx='1'/>` +
		`<circle cx='62' cy='22' r='1.4' fill='var(--accent)' stroke='none'/>` +
		`<circle cx='62' cy='38' r='1.4' fill='var(--accent)' stroke='none'/>` +
		`<line x1='35' y1='22' x2='55' y2='22' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='35' y1='38' x2='55' y2='38' stroke='var(--ink)' stroke-width='1'/>` +
		`</g>`,

	// cache: lightning in a box
	"cache": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='32' y='14' width='36' height='32' rx='2'/>` +
		`<path d='M 52 18 L 42 32 L 50 32 L 46 42 L 58 26 L 50 26 Z' fill='var(--accent)' stroke='var(--accent)' stroke-width='1.2' stroke-linejoin='round'/>` +
		`</g>`,

	// stream: parallel flowing lines
	"stream": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<path d='M 28 22 Q 40 16 50 22 T 72 22'/>` +
		`<path d='M 28 30 Q 40 24 50 30 T 72 30' stroke='var(--ink)' stroke-width='1.2'/>` +
		`<path d='M 28 38 Q 40 32 50 38 T 72 38'/>` +
		`<polygon points='72,22 68,20 68,24' fill='var(--accent)' stroke='none'/>` +
		`<polygon points='72,38 68,36 68,40' fill='var(--accent)' stroke='none'/>` +
		`</g>`,

	// ── Network ─────────────────────────────────────────────────────
	// load-balancer: split/fanout diamond
	"load-balancer": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<polygon points='50,14 64,30 50,46 36,30' fill='var(--accent-soft)'/>` +
		`<line x1='50' y1='30' x2='50' y2='22' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='50' y1='30' x2='42' y2='34' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='50' y1='30' x2='58' y2='34' stroke='var(--ink)' stroke-width='1'/>` +
		`<polygon points='42,34 44,32 45,36' fill='var(--ink)' stroke='none'/>` +
		`<polygon points='58,34 55,36 56,32' fill='var(--ink)' stroke='none'/>` +
		`</g>`,

	// gateway: opening between two pillars
	"gateway": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='30' y='14' width='8' height='32'/>` +
		`<rect x='62' y='14' width='8' height='32'/>` +
		`<path d='M 38 18 L 62 18' stroke='var(--ink)' stroke-width='1.2'/>` +
		`<path d='M 42 30 L 58 30' stroke='var(--accent)' stroke-width='1.6'/>` +
		`<polygon points='58,30 54,28 54,32' fill='var(--accent)' stroke='none'/>` +
		`</g>`,

	// cdn: globe with edge nodes
	"cdn": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<circle cx='50' cy='30' r='12'/>` +
		`<ellipse cx='50' cy='30' rx='12' ry='4' stroke='var(--ink)' stroke-width='1'/>` +
		`<ellipse cx='50' cy='30' rx='5' ry='12' stroke='var(--ink)' stroke-width='1'/>` +
		`<circle cx='32' cy='18' r='2.2' fill='var(--accent)' stroke='none'/>` +
		`<circle cx='68' cy='18' r='2.2' fill='var(--accent)' stroke='none'/>` +
		`<circle cx='32' cy='42' r='2.2' fill='var(--accent)' stroke='none'/>` +
		`<circle cx='68' cy='42' r='2.2' fill='var(--accent)' stroke='none'/>` +
		`</g>`,

	// dns: tag-style label with dot
	"dns": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<path d='M 30 22 L 56 22 L 66 30 L 56 38 L 30 38 Z'/>` +
		`<circle cx='38' cy='30' r='2' fill='var(--accent)' stroke='none'/>` +
		`<text x='44' y='33' font-family='IBM Plex Mono' font-size='7' fill='var(--ink)'>.io</text>` +
		`</g>`,

	// vpc: dashed cloud-ish container
	"vpc": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round' stroke-dasharray='3 2'>` +
		`<rect x='26' y='14' width='48' height='32' rx='6'/>` +
		`</g>` +
		`<g fill='none' stroke='var(--ink)' stroke-width='1'>` +
		`<rect x='34' y='22' width='10' height='8'/>` +
		`<rect x='48' y='22' width='10' height='8'/>` +
		`<rect x='42' y='34' width='10' height='8'/>` +
		`</g>`,

	// subnet: nested rect inside dashed outer
	"subnet": `<g fill='none' stroke='var(--ink)' stroke-width='1' stroke-dasharray='2 2'>` +
		`<rect x='28' y='16' width='44' height='28' rx='2'/>` +
		`</g>` +
		`<g fill='var(--accent-soft)' stroke='var(--accent)' stroke-width='1.5'>` +
		`<rect x='38' y='24' width='24' height='12' rx='1.5'/>` +
		`</g>`,

	// firewall: brick wall + flame
	"firewall": `<g fill='none' stroke='var(--accent)' stroke-width='1.4' stroke-linejoin='round'>` +
		`<rect x='30' y='24' width='40' height='20'/>` +
		`<line x1='30' y1='34' x2='70' y2='34' stroke='var(--accent)' stroke-width='1.2'/>` +
		`<line x1='40' y1='24' x2='40' y2='34' stroke='var(--accent)' stroke-width='1.2'/>` +
		`<line x1='52' y1='24' x2='52' y2='34' stroke='var(--accent)' stroke-width='1.2'/>` +
		`<line x1='62' y1='24' x2='62' y2='34' stroke='var(--accent)' stroke-width='1.2'/>` +
		`<line x1='36' y1='34' x2='36' y2='44' stroke='var(--accent)' stroke-width='1.2'/>` +
		`<line x1='48' y1='34' x2='48' y2='44' stroke='var(--accent)' stroke-width='1.2'/>` +
		`<line x1='58' y1='34' x2='58' y2='44' stroke='var(--accent)' stroke-width='1.2'/>` +
		`<path d='M 50 22 Q 44 18 48 14 Q 54 16 52 12 Q 58 16 54 22 Z' fill='var(--accent)' stroke='var(--accent)' stroke-width='1'/>` +
		`</g>`,

	// ── Services ────────────────────────────────────────────────────
	// api: square brackets around dots ({ … })
	"api": `<g fill='none' stroke='var(--accent)' stroke-width='1.6' stroke-linejoin='round' stroke-linecap='round'>` +
		`<path d='M 38 16 Q 30 16 30 22 L 30 28 Q 30 30 28 30 Q 30 30 30 32 L 30 38 Q 30 44 38 44'/>` +
		`<path d='M 62 16 Q 70 16 70 22 L 70 28 Q 70 30 72 30 Q 70 30 70 32 L 70 38 Q 70 44 62 44'/>` +
		`<circle cx='44' cy='30' r='1.8' fill='var(--accent)' stroke='none'/>` +
		`<circle cx='50' cy='30' r='1.8' fill='var(--accent)' stroke='none'/>` +
		`<circle cx='56' cy='30' r='1.8' fill='var(--accent)' stroke='none'/>` +
		`</g>`,

	// queue: stacked envelopes / message line
	"queue": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='32' y='22' width='10' height='16' fill='var(--accent-soft)'/>` +
		`<rect x='44' y='22' width='10' height='16' fill='var(--accent-soft)'/>` +
		`<rect x='56' y='22' width='10' height='16' fill='var(--accent-soft)'/>` +
		`<line x1='66' y1='30' x2='72' y2='30' stroke='var(--accent)' stroke-width='1.4'/>` +
		`<polygon points='72,30 68,28 68,32' fill='var(--accent)' stroke='none'/>` +
		`<line x1='28' y1='30' x2='32' y2='30' stroke='var(--ink)' stroke-width='1'/>` +
		`</g>`,

	// topic: pub/sub fanout from a single node
	"topic": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<circle cx='34' cy='30' r='5' fill='var(--accent)' stroke='var(--accent)'/>` +
		`<g stroke='var(--ink)' stroke-width='1'>` +
		`<line x1='39' y1='30' x2='60' y2='18'/>` +
		`<line x1='39' y1='30' x2='60' y2='30'/>` +
		`<line x1='39' y1='30' x2='60' y2='42'/>` +
		`</g>` +
		`<circle cx='64' cy='18' r='3'/>` +
		`<circle cx='64' cy='30' r='3'/>` +
		`<circle cx='64' cy='42' r='3'/>` +
		`</g>`,

	// ── Observability ───────────────────────────────────────────────
	// log: lined document
	"log": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<path d='M 34 14 L 60 14 L 68 22 L 68 46 L 34 46 Z'/>` +
		`<path d='M 60 14 L 60 22 L 68 22' stroke='var(--accent)' stroke-width='1.4'/>` +
		`<line x1='38' y1='28' x2='62' y2='28' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='38' y1='34' x2='62' y2='34' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='38' y1='40' x2='54' y2='40' stroke='var(--ink)' stroke-width='1'/>` +
		`</g>`,

	// metric: small bar chart
	"metric": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='30' y='14' width='40' height='32' fill='none' stroke='var(--ink)' stroke-width='1'/>` +
		`<rect x='36' y='34' width='5' height='8' fill='var(--accent)' stroke='none'/>` +
		`<rect x='44' y='28' width='5' height='14' fill='var(--accent)' stroke='none'/>` +
		`<rect x='52' y='22' width='5' height='20' fill='var(--accent)' stroke='none'/>` +
		`<rect x='60' y='18' width='5' height='24' fill='var(--accent)' stroke='none'/>` +
		`</g>`,

	// trace: waterfall spans
	"trace": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='30' y='18' width='36' height='3' fill='var(--accent)' stroke='none'/>` +
		`<rect x='34' y='24' width='24' height='3' fill='var(--accent)' opacity='0.75' stroke='none'/>` +
		`<rect x='38' y='30' width='14' height='3' fill='var(--accent)' opacity='0.55' stroke='none'/>` +
		`<rect x='42' y='36' width='10' height='3' fill='var(--accent)' opacity='0.4' stroke='none'/>` +
		`<line x1='28' y1='14' x2='28' y2='44' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='28' y1='44' x2='72' y2='44' stroke='var(--ink)' stroke-width='1'/>` +
		`</g>`,

	// ── Identity ────────────────────────────────────────────────────
	// user: head + shoulders
	"user": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<circle cx='50' cy='22' r='7'/>` +
		`<path d='M 34 46 Q 34 32 50 32 Q 66 32 66 46'/>` +
		`</g>`,

	// mobile: phone outline + speaker
	"mobile": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='40' y='12' width='20' height='36' rx='3'/>` +
		`<line x1='46' y1='16' x2='54' y2='16' stroke='var(--ink)' stroke-width='1'/>` +
		`<circle cx='50' cy='44' r='1.2' fill='var(--ink)' stroke='none'/>` +
		`</g>`,

	// browser: window with tab + url bar
	"browser": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='28' y='14' width='44' height='32' rx='2'/>` +
		`<line x1='28' y1='22' x2='72' y2='22' stroke='var(--accent)' stroke-width='1.4'/>` +
		`<circle cx='32' cy='18' r='1.2' fill='var(--ink)' stroke='none'/>` +
		`<circle cx='36' cy='18' r='1.2' fill='var(--ink)' stroke='none'/>` +
		`<circle cx='40' cy='18' r='1.2' fill='var(--ink)' stroke='none'/>` +
		`<rect x='34' y='28' width='32' height='4' fill='var(--accent-soft)' stroke='var(--ink)' stroke-width='0.8'/>` +
		`</g>`,

	// ── Security ────────────────────────────────────────────────────
	// secret: padlock with asterisks inside
	"secret": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<rect x='38' y='28' width='24' height='18' rx='2' fill='var(--accent)' stroke='var(--accent)'/>` +
		`<path d='M 42 28 V 22 a 8 8 0 0 1 16 0 V 28'/>` +
		`<text x='50' y='42' text-anchor='middle' font-family='IBM Plex Mono' font-size='10' font-weight='700' fill='var(--bg)'>***</text>` +
		`</g>`,

	// key-icon: classic key shape
	"key-icon": `<g fill='none' stroke='var(--accent)' stroke-width='1.5' stroke-linejoin='round'>` +
		`<circle cx='36' cy='30' r='8' fill='var(--accent-soft)'/>` +
		`<circle cx='36' cy='30' r='2.4' fill='var(--accent)' stroke='none'/>` +
		`<line x1='44' y1='30' x2='70' y2='30'/>` +
		`<line x1='62' y1='30' x2='62' y2='36' stroke='var(--accent)' stroke-width='1.6'/>` +
		`<line x1='68' y1='30' x2='68' y2='38' stroke='var(--accent)' stroke-width='1.6'/>` +
		`</g>`,
}

// ArchIcon returns the SVG `<g>` snippet for one named arch icon, or "" if unknown.
// Used by the box: DSL when an option references ?icon=<name>, and by the
// arch:* dispatcher branch (which wraps the snippet in a sized <svg>).
func ArchIcon(name string) string {
	if g, ok := archLibrary[name]; ok {
		return g
	}
	return ""
}

// ArchLibrary exposes the underlying registry as a read-only-by-convention
// map. Tests iterate over every entry; external callers should use
// ArchIcon for lookups.
var ArchLibrary = archLibrary
