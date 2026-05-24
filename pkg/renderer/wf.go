package renderer

// WFLibrary — named wireframe registry. Agents reference these via
// `"art": "wf:<name>"`. Each entry is a self-contained SVG fragment
// using theme CSS variables (var(--accent), var(--ink), etc.) so it
// inherits the active theme automatically.
//
// Convention: every entry uses viewBox='0 0 100 60' and a top-level
// `style='width:80%;height:auto'` so they slot into the 4:3 option-art
// area at consistent visual size.
//
// Categories (50 entries total):
//   calendars/dates (7), charts (7), identities (5), phone screens (7),
//   controls (6), status (5), typography (3), connections (4),
//   metaphors (3), misc (3).
var WFLibrary = map[string]string{

	// ── calendars / dates ────────────────────────────────────────────
	"calendar-grid": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'><g>` +
		gridCells(7, 5, 14, 6, 11, 7, []int{0, 2, 5, 7, 10, 12, 15, 18, 22, 25, 27}) + `</g></svg>`,

	"calendar-week": `<svg viewBox='0 0 100 60' style='width:82%;height:auto'>` +
		`<g font-family='IBM Plex Mono' font-size='6' fill='var(--ink-3)' letter-spacing='0.04em'><text x='9' y='10'>MON</text><text x='22' y='10'>TUE</text><text x='35' y='10'>WED</text><text x='48' y='10'>THU</text><text x='61' y='10'>FRI</text><text x='74' y='10'>SAT</text><text x='87' y='10'>SUN</text></g>` +
		`<rect x='6' y='14' width='12' height='28' fill='var(--accent)' opacity='0.7'/><rect x='19' y='14' width='12' height='28' fill='var(--accent)'/><rect x='32' y='14' width='12' height='28' fill='var(--accent)' opacity='0.3'/><rect x='45' y='14' width='12' height='28' fill='var(--accent)' opacity='0.85'/><rect x='58' y='14' width='12' height='28' fill='none' stroke='var(--rule-2)' stroke-width='0.6'/><rect x='71' y='14' width='12' height='28' fill='var(--accent)' opacity='0.5'/><rect x='84' y='14' width='12' height='28' fill='none' stroke='var(--rule-2)' stroke-width='0.6'/></svg>`,

	"streak-counter": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<text x='50' y='36' text-anchor='middle' font-family='IBM Plex Mono' font-size='30' font-weight='700' fill='var(--accent)' letter-spacing='-0.04em'>27</text>` +
		`<text x='50' y='46' text-anchor='middle' font-family='IBM Plex Mono' font-size='5' fill='var(--ink-3)' letter-spacing='0.12em'>DAY STREAK</text>` +
		`<g transform='translate(34,50)'><circle cx='0' cy='0' r='2.5' fill='var(--accent)'/><circle cx='6' cy='0' r='2.5' fill='var(--accent)'/><circle cx='12' cy='0' r='2.5' fill='var(--accent)'/><circle cx='18' cy='0' r='2.5' fill='var(--accent)'/><circle cx='24' cy='0' r='2.5' fill='var(--accent)'/><circle cx='30' cy='0' r='2.5' fill='var(--accent)'/><circle cx='36' cy='0' r='2.5' fill='none' stroke='var(--ink)' stroke-width='0.8'/></g></svg>`,

	"day-strip": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<g><circle cx='14' cy='30' r='6' fill='var(--accent)'/><text x='14' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' fill='var(--bg)' font-weight='700'>M</text>` +
		`<circle cx='28' cy='30' r='6' fill='var(--accent)'/><text x='28' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' fill='var(--bg)' font-weight='700'>T</text>` +
		`<circle cx='42' cy='30' r='6' fill='var(--accent)'/><text x='42' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' fill='var(--bg)' font-weight='700'>W</text>` +
		`<circle cx='56' cy='30' r='6' fill='var(--accent)'/><text x='56' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' fill='var(--bg)' font-weight='700'>T</text>` +
		`<circle cx='70' cy='30' r='6' fill='var(--accent)'/><text x='70' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' fill='var(--bg)' font-weight='700'>F</text>` +
		`<circle cx='84' cy='30' r='6' fill='none' stroke='var(--ink)' stroke-width='1'/><text x='84' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' fill='var(--ink-3)' font-weight='700'>S</text></g></svg>`,

	"year-heatmap": `<svg viewBox='0 0 100 60' style='width:84%;height:auto'><g>` +
		yearGrid() + `</g></svg>`,

	"time-of-day": `<svg viewBox='0 0 100 60' style='width:78%;height:auto'>` +
		`<g transform='translate(20,30)'><circle cx='0' cy='0' r='8' fill='var(--accent)'/><g stroke='var(--accent)' stroke-width='1.4' stroke-linecap='round'><line x1='-13' y1='0' x2='-10' y2='0'/><line x1='10' y1='0' x2='13' y2='0'/><line x1='0' y1='-13' x2='0' y2='-10'/><line x1='0' y1='10' x2='0' y2='13'/></g></g>` +
		`<g transform='translate(50,30)'><circle cx='0' cy='0' r='8' fill='none' stroke='var(--ink-2)' stroke-width='1.2'/><path d='M 0 -8 A 8 8 0 0 1 0 8 Z' fill='var(--ink-2)'/></g>` +
		`<g transform='translate(80,30)'><path d='M -5 -7 A 7 7 0 1 0 5 7 A 8 8 0 0 1 -5 -7 Z' fill='var(--ink)'/></g></svg>`,

	"clock": `<svg viewBox='0 0 100 60' style='width:70%;height:auto'>` +
		`<circle cx='50' cy='30' r='22' fill='none' stroke='var(--ink)' stroke-width='1.6'/>` +
		`<g stroke='var(--ink-3)' stroke-width='0.8'><line x1='50' y1='10' x2='50' y2='13'/><line x1='70' y1='30' x2='67' y2='30'/><line x1='50' y1='50' x2='50' y2='47'/><line x1='30' y1='30' x2='33' y2='30'/></g>` +
		`<line x1='50' y1='30' x2='50' y2='17' stroke='var(--accent)' stroke-width='1.8' stroke-linecap='round'/>` +
		`<line x1='50' y1='30' x2='62' y2='34' stroke='var(--accent)' stroke-width='1.4' stroke-linecap='round'/>` +
		`<circle cx='50' cy='30' r='1.8' fill='var(--accent)'/></svg>`,

	// ── charts ───────────────────────────────────────────────────────
	"spark-rising": `<svg viewBox='0 0 100 60' style='width:82%;height:auto'>` +
		`<path d='M 10 48 L 22 44 L 34 40 L 46 35 L 58 28 L 70 22 L 82 14 L 90 8 L 90 52 L 10 52 Z' fill='var(--accent)' opacity='0.15'/>` +
		`<path d='M 10 48 L 22 44 L 34 40 L 46 35 L 58 28 L 70 22 L 82 14 L 90 8' fill='none' stroke='var(--accent)' stroke-width='1.6' stroke-linejoin='round'/>` +
		`<circle cx='90' cy='8' r='2.2' fill='var(--accent)'/></svg>`,

	"spark-flat": `<svg viewBox='0 0 100 60' style='width:82%;height:auto'>` +
		`<path d='M 10 32 L 22 30 L 34 33 L 46 29 L 58 32 L 70 30 L 82 31 L 90 30 L 90 52 L 10 52 Z' fill='var(--accent)' opacity='0.15'/>` +
		`<path d='M 10 32 L 22 30 L 34 33 L 46 29 L 58 32 L 70 30 L 82 31 L 90 30' fill='none' stroke='var(--accent)' stroke-width='1.6' stroke-linejoin='round'/></svg>`,

	"spark-noisy": `<svg viewBox='0 0 100 60' style='width:82%;height:auto'>` +
		`<path d='M 10 40 L 18 18 L 26 35 L 34 12 L 42 38 L 50 22 L 58 42 L 66 16 L 74 32 L 82 20 L 90 36' fill='none' stroke='var(--accent)' stroke-width='1.3' stroke-linejoin='round'/></svg>`,

	"bars-up": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<rect x='12' y='42' width='8' height='10' fill='var(--accent)' opacity='0.7'/>` +
		`<rect x='24' y='36' width='8' height='16' fill='var(--accent)' opacity='0.75'/>` +
		`<rect x='36' y='30' width='8' height='22' fill='var(--accent)' opacity='0.8'/>` +
		`<rect x='48' y='32' width='8' height='20' fill='var(--accent)' opacity='0.8'/>` +
		`<rect x='60' y='22' width='8' height='30' fill='var(--accent)' opacity='0.9'/>` +
		`<rect x='72' y='14' width='8' height='38' fill='var(--accent)'/>` +
		`<rect x='84' y='8' width='8' height='44' fill='var(--accent)'/></svg>`,

	"donut": `<svg viewBox='0 0 100 60' style='width:70%;height:auto'>` +
		`<circle cx='50' cy='30' r='20' fill='none' stroke='var(--ink-3)' stroke-width='8' opacity='0.3'/>` +
		`<circle cx='50' cy='30' r='20' fill='none' stroke='var(--accent)' stroke-width='8' stroke-dasharray='80 126' stroke-linecap='butt' transform='rotate(-90 50 30)'/>` +
		`<text x='50' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='10' font-weight='700' fill='var(--ink)'>64%</text></svg>`,

	"gauge": `<svg viewBox='0 0 100 60' style='width:78%;height:auto'>` +
		`<path d='M 18 48 A 32 32 0 0 1 82 48' fill='none' stroke='var(--ink-3)' stroke-width='6' opacity='0.3'/>` +
		`<path d='M 18 48 A 32 32 0 0 1 70 24' fill='none' stroke='var(--accent)' stroke-width='6'/>` +
		`<circle cx='70' cy='24' r='3.5' fill='var(--accent)'/>` +
		`<text x='50' y='44' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' font-weight='600' fill='var(--ink)'>72%</text></svg>`,

	"dot-trend": `<svg viewBox='0 0 100 60' style='width:84%;height:auto'>` +
		`<g><circle cx='15' cy='45' r='2' fill='var(--ink-3)'/><circle cx='27' cy='40' r='2' fill='var(--ink-3)'/><circle cx='39' cy='36' r='2' fill='var(--ink-3)'/><circle cx='51' cy='30' r='2' fill='var(--ink-3)'/><circle cx='63' cy='22' r='2' fill='var(--accent)'/><circle cx='75' cy='18' r='2' fill='var(--accent)'/><circle cx='87' cy='12' r='2.5' fill='var(--accent)'/></g>` +
		`<line x1='15' y1='45' x2='87' y2='12' stroke='var(--accent)' stroke-width='0.6' stroke-dasharray='2 2' opacity='0.5'/></svg>`,

	// ── identities / avatars ─────────────────────────────────────────
	"avatar-single": `<svg viewBox='0 0 100 60' style='width:60%;height:auto'>` +
		`<circle cx='50' cy='28' r='14' fill='var(--surface)' stroke='var(--ink)' stroke-width='1.4'/>` +
		`<text x='50' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='13' font-weight='600' fill='var(--ink)'>M</text>` +
		`<text x='50' y='52' text-anchor='middle' font-family='IBM Plex Mono' font-size='5' fill='var(--ink-3)' letter-spacing='0.1em'>· JUST YOU ·</text></svg>`,

	"avatar-pair": `<svg viewBox='0 0 100 60' style='width:60%;height:auto'>` +
		`<circle cx='40' cy='30' r='12' fill='var(--surface)' stroke='var(--ink)' stroke-width='1.2'/>` +
		`<text x='40' y='34' text-anchor='middle' font-family='IBM Plex Mono' font-size='10' font-weight='600' fill='var(--ink)'>M</text>` +
		`<circle cx='60' cy='30' r='12' fill='var(--accent)' stroke='var(--accent)' stroke-width='1.2'/>` +
		`<text x='60' y='34' text-anchor='middle' font-family='IBM Plex Mono' font-size='10' font-weight='600' fill='var(--bg)'>J</text></svg>`,

	"avatar-circle": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<circle cx='50' cy='30' r='22' fill='none' stroke='var(--rule-2)' stroke-width='0.6' stroke-dasharray='2 2'/>` +
		`<circle cx='50' cy='30' r='9' fill='var(--accent)' stroke='var(--accent)' stroke-width='1'/><text x='50' y='34' text-anchor='middle' font-family='IBM Plex Mono' font-size='8' font-weight='600' fill='var(--bg)'>M</text>` +
		`<circle cx='28' cy='14' r='6' fill='var(--surface)' stroke='var(--ink)' stroke-width='1'/><text x='28' y='17' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' font-weight='600' fill='var(--ink)'>J</text>` +
		`<circle cx='72' cy='14' r='6' fill='var(--surface)' stroke='var(--ink)' stroke-width='1'/><text x='72' y='17' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' font-weight='600' fill='var(--ink)'>S</text>` +
		`<circle cx='22' cy='44' r='6' fill='var(--surface)' stroke='var(--ink)' stroke-width='1'/><text x='22' y='47' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' font-weight='600' fill='var(--ink)'>A</text>` +
		`<circle cx='78' cy='44' r='6' fill='var(--surface)' stroke='var(--ink)' stroke-width='1'/><text x='78' y='47' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' font-weight='600' fill='var(--ink)'>K</text></svg>`,

	"avatar-feed": `<svg viewBox='0 0 100 60' style='width:84%;height:auto'><g font-family='IBM Plex Mono' font-size='6' font-weight='600'>` +
		`<circle cx='12' cy='12' r='5' fill='var(--surface)' stroke='var(--ink)' stroke-width='0.8'/><text x='12' y='14.5' text-anchor='middle' fill='var(--ink)'>J</text><line x1='22' y1='12' x2='75' y2='12' stroke='var(--ink)' stroke-width='3' opacity='0.45'/><text x='86' y='14' fill='var(--accent)'>♥ 12</text>` +
		`<line x1='4' y1='22' x2='96' y2='22' stroke='var(--rule)' stroke-width='0.4'/>` +
		`<circle cx='12' cy='32' r='5' fill='var(--surface)' stroke='var(--ink)' stroke-width='0.8'/><text x='12' y='34.5' text-anchor='middle' fill='var(--ink)'>S</text><line x1='22' y1='32' x2='70' y2='32' stroke='var(--ink)' stroke-width='3' opacity='0.45'/><text x='86' y='34' fill='var(--ink-3)'>♥ 4</text>` +
		`<line x1='4' y1='42' x2='96' y2='42' stroke='var(--rule)' stroke-width='0.4'/>` +
		`<circle cx='12' cy='52' r='5' fill='var(--surface)' stroke='var(--ink)' stroke-width='0.8'/><text x='12' y='54.5' text-anchor='middle' fill='var(--ink)'>M</text><line x1='22' y1='52' x2='65' y2='52' stroke='var(--ink)' stroke-width='3' opacity='0.45'/><text x='86' y='54' fill='var(--ink-3)'>♥ 1</text></g></svg>`,

	"avatar-private": `<svg viewBox='0 0 100 60' style='width:60%;height:auto'>` +
		`<circle cx='50' cy='28' r='14' fill='var(--surface)' stroke='var(--ink)' stroke-width='1.4'/>` +
		`<text x='50' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='13' font-weight='600' fill='var(--ink)'>M</text>` +
		`<g transform='translate(60,38)'><rect x='-4' y='-2' width='8' height='6' rx='1' fill='var(--accent)'/><path d='M -2.5 -2 V -4 a 2.5 2.5 0 0 1 5 0 V -2' fill='none' stroke='var(--accent)' stroke-width='1'/></g></svg>`,

	// ── phone screens ────────────────────────────────────────────────
	"phone-blank":   phoneScreen(``),
	"phone-list":    phoneScreen(`<rect x='4' y='8' width='34' height='3' fill='var(--ink)' opacity='0.6'/><line x1='2' y1='15' x2='40' y2='15' stroke='var(--rule)' stroke-width='0.3'/><rect x='4' y='18' width='28' height='3' fill='var(--ink)' opacity='0.6'/><line x1='2' y1='25' x2='40' y2='25' stroke='var(--rule)' stroke-width='0.3'/><rect x='4' y='28' width='32' height='3' fill='var(--ink)' opacity='0.6'/><line x1='2' y1='35' x2='40' y2='35' stroke='var(--rule)' stroke-width='0.3'/><rect x='4' y='38' width='24' height='3' fill='var(--ink)' opacity='0.6'/>`),
	"phone-card":    phoneScreen(`<rect x='4' y='8' width='34' height='28' fill='none' stroke='var(--accent)' stroke-width='0.8'/><rect x='6' y='10' width='14' height='14' fill='var(--accent)' opacity='0.3'/><rect x='6' y='28' width='22' height='2.5' fill='var(--ink)' opacity='0.6'/><rect x='6' y='32' width='18' height='2.5' fill='var(--ink)' opacity='0.4'/>`),
	"phone-input":   phoneScreen(`<rect x='4' y='18' width='34' height='8' fill='var(--surface)' stroke='var(--ink)' stroke-width='0.8'/><line x1='32' y1='20' x2='32' y2='24' stroke='var(--ink)' stroke-width='0.6'/><rect x='4' y='30' width='34' height='6' fill='var(--accent)'/>`),
	"phone-tabs":    phoneScreen(`<rect x='4' y='8' width='34' height='28' fill='none' stroke='var(--rule)' stroke-width='0.4'/><line x1='2' y1='44' x2='40' y2='44' stroke='var(--rule-2)' stroke-width='0.4'/><circle cx='10' cy='49' r='1.5' fill='var(--accent)'/><circle cx='21' cy='49' r='1.5' fill='var(--ink-3)' opacity='0.5'/><circle cx='32' cy='49' r='1.5' fill='var(--ink-3)' opacity='0.5'/>`),
	"phone-onboard": phoneScreen(`<rect x='4' y='8' width='14' height='2' fill='var(--accent)'/><rect x='18' y='8' width='20' height='2' fill='var(--rule-2)'/><rect x='4' y='16' width='28' height='4' fill='var(--ink)' opacity='0.7'/><rect x='4' y='22' width='18' height='4' fill='var(--ink)' opacity='0.7'/><rect x='4' y='32' width='34' height='5' fill='var(--accent)'/><rect x='4' y='40' width='34' height='5' fill='none' stroke='var(--ink)' stroke-width='0.6'/>`),
	"phone-stats":   phoneScreen(`<rect x='4' y='8' width='12' height='2.5' fill='var(--ink)' opacity='0.6'/><text x='4' y='20' font-family='IBM Plex Mono' font-size='10' font-weight='700' fill='var(--accent)'>27</text><rect x='4' y='28' width='34' height='14' fill='none' stroke='var(--rule-2)' stroke-width='0.4'/><path d='M 5 40 L 11 36 L 17 38 L 23 32 L 29 30 L 35 26' fill='none' stroke='var(--accent)' stroke-width='1.2'/>`),

	// ── controls ─────────────────────────────────────────────────────
	"toggle-on": `<svg viewBox='0 0 100 60' style='width:55%;height:auto'>` +
		`<rect x='32' y='22' width='36' height='18' rx='9' fill='var(--accent)' stroke='var(--accent)' stroke-width='1'/>` +
		`<circle cx='59' cy='31' r='6.5' fill='var(--bg)'/>` +
		`<text x='50' y='55' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' fill='var(--accent)' letter-spacing='0.1em' font-weight='600'>ON</text></svg>`,

	"toggle-off": `<svg viewBox='0 0 100 60' style='width:55%;height:auto'>` +
		`<rect x='32' y='22' width='36' height='18' rx='9' fill='none' stroke='var(--ink)' stroke-width='1.4'/>` +
		`<circle cx='41' cy='31' r='6.5' fill='var(--ink)'/>` +
		`<text x='50' y='55' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' fill='var(--ink-3)' letter-spacing='0.1em'>OFF</text></svg>`,

	"button-accent": `<svg viewBox='0 0 100 60' style='width:70%;height:auto'>` +
		`<rect x='14' y='22' width='72' height='18' fill='var(--accent)' stroke='var(--ink)' stroke-width='1'/>` +
		`<text x='50' y='34' text-anchor='middle' font-family='IBM Plex Sans' font-size='8' font-weight='700' fill='var(--bg)' letter-spacing='0.04em'>CONTINUE</text></svg>`,

	"button-ghost": `<svg viewBox='0 0 100 60' style='width:70%;height:auto'>` +
		`<rect x='14' y='22' width='72' height='18' fill='none' stroke='var(--ink)' stroke-width='1' stroke-dasharray='3 2'/>` +
		`<text x='50' y='34' text-anchor='middle' font-family='IBM Plex Sans' font-size='8' font-weight='600' fill='var(--ink-2)' letter-spacing='0.04em'>maybe later</text></svg>`,

	"slider": `<svg viewBox='0 0 100 60' style='width:78%;height:auto'>` +
		`<line x1='14' y1='30' x2='86' y2='30' stroke='var(--ink)' stroke-width='1.5' opacity='0.3'/>` +
		`<line x1='14' y1='30' x2='60' y2='30' stroke='var(--accent)' stroke-width='2'/>` +
		`<circle cx='60' cy='30' r='5' fill='var(--accent)' stroke='var(--ink)' stroke-width='1'/>` +
		`<text x='86' y='46' text-anchor='end' font-family='IBM Plex Mono' font-size='8' font-weight='700' fill='var(--accent)'>72</text></svg>`,

	"dropdown": `<svg viewBox='0 0 100 60' style='width:78%;height:auto'>` +
		`<rect x='14' y='22' width='72' height='14' fill='var(--surface)' stroke='var(--ink)' stroke-width='1'/>` +
		`<text x='18' y='32' font-family='IBM Plex Mono' font-size='7' fill='var(--ink-2)'>select option</text>` +
		`<path d='M 78 27 L 82 31 L 78 35' fill='none' stroke='var(--ink)' stroke-width='1.2' stroke-linecap='round'/></svg>`,

	// ── status ───────────────────────────────────────────────────────
	"badge-new": `<svg viewBox='0 0 100 60' style='width:60%;height:auto'>` +
		`<rect x='28' y='20' width='44' height='20' fill='var(--accent)'/>` +
		`<text x='50' y='34' text-anchor='middle' font-family='IBM Plex Mono' font-size='10' font-weight='700' fill='var(--bg)' letter-spacing='0.16em'>NEW</text></svg>`,

	"pill-row": `<svg viewBox='0 0 100 60' style='width:84%;height:auto'>` +
		`<rect x='8' y='22' width='22' height='14' rx='7' fill='var(--accent)' stroke='var(--accent)' stroke-width='0.8'/><text x='19' y='31.5' text-anchor='middle' font-family='IBM Plex Mono' font-size='6.5' fill='var(--bg)' font-weight='600'>active</text>` +
		`<rect x='34' y='22' width='28' height='14' rx='7' fill='none' stroke='var(--ink)' stroke-width='0.8'/><text x='48' y='31.5' text-anchor='middle' font-family='IBM Plex Mono' font-size='6.5' fill='var(--ink)' font-weight='600'>pending</text>` +
		`<rect x='66' y='22' width='26' height='14' rx='7' fill='none' stroke='var(--ink-3)' stroke-width='0.8' stroke-dasharray='2 1'/><text x='79' y='31.5' text-anchor='middle' font-family='IBM Plex Mono' font-size='6.5' fill='var(--ink-3)' font-weight='600'>draft</text></svg>`,

	"snowflake": `<svg viewBox='0 0 100 60' style='width:55%;height:auto'>` +
		`<g stroke='var(--accent)' stroke-width='1.8' fill='none' stroke-linecap='round' transform='translate(50 30)'><line x1='0' y1='-22' x2='0' y2='22'/><line x1='-19' y1='-11' x2='19' y2='11'/><line x1='-19' y1='11' x2='19' y2='-11'/><line x1='-4' y1='-19' x2='0' y2='-22'/><line x1='4' y1='-19' x2='0' y2='-22'/><line x1='-4' y1='19' x2='0' y2='22'/><line x1='4' y1='19' x2='0' y2='22'/><line x1='-16' y1='-14' x2='-19' y2='-11'/><line x1='-16' y1='-8' x2='-19' y2='-11'/></g></svg>`,

	"lock": `<svg viewBox='0 0 100 60' style='width:45%;height:auto'>` +
		`<rect x='38' y='28' width='24' height='18' rx='2' fill='var(--accent)'/>` +
		`<path d='M 42 28 V 22 a 8 8 0 0 1 16 0 V 28' fill='none' stroke='var(--accent)' stroke-width='2'/>` +
		`<circle cx='50' cy='37' r='2.5' fill='var(--bg)'/></svg>`,

	"check-large": `<svg viewBox='0 0 100 60' style='width:50%;height:auto'>` +
		`<circle cx='50' cy='30' r='20' fill='var(--accent)'/>` +
		`<path d='M 38 30 L 47 38 L 62 22' fill='none' stroke='var(--bg)' stroke-width='3.5' stroke-linecap='round' stroke-linejoin='round'/></svg>`,

	// ── typography samples ───────────────────────────────────────────
	"serif-sample": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<text x='50' y='38' text-anchor='middle' font-family='IBM Plex Serif' font-style='italic' font-size='22' font-weight='500' fill='var(--ink)'>Tide</text>` +
		`<text x='50' y='50' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' fill='var(--ink-3)' letter-spacing='0.1em'>EDITORIAL · SERIF</text></svg>`,

	"sans-sample": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<text x='50' y='37' text-anchor='middle' font-family='IBM Plex Sans' font-weight='700' font-size='22' fill='var(--ink)' letter-spacing='-0.02em'>Tide</text>` +
		`<text x='50' y='50' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' fill='var(--ink-3)' letter-spacing='0.1em'>MODERN · SANS</text></svg>`,

	"mono-sample": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<text x='50' y='37' text-anchor='middle' font-family='IBM Plex Mono' font-weight='500' font-size='20' fill='var(--ink)'>tide_</text>` +
		`<text x='50' y='50' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' fill='var(--ink-3)' letter-spacing='0.1em'>TECHNICAL · MONO</text></svg>`,

	// ── connections / graphs ─────────────────────────────────────────
	"graph-force": `<svg viewBox='0 0 100 60' style='width:84%;height:auto'>` +
		`<g stroke='var(--ink-3)' stroke-width='0.6' opacity='0.6'><line x1='50' y1='30' x2='25' y2='15'/><line x1='50' y1='30' x2='75' y2='17'/><line x1='50' y1='30' x2='42' y2='52'/><line x1='50' y1='30' x2='72' y2='48'/><line x1='25' y1='15' x2='10' y2='28'/><line x1='75' y1='17' x2='88' y2='32'/><line x1='42' y1='52' x2='72' y2='48'/></g>` +
		`<g><circle cx='50' cy='30' r='5' fill='var(--accent)'/><circle cx='25' cy='15' r='3.5' fill='var(--ink)' opacity='0.7'/><circle cx='75' cy='17' r='3.5' fill='var(--ink)' opacity='0.7'/><circle cx='42' cy='52' r='3' fill='var(--ink)' opacity='0.6'/><circle cx='72' cy='48' r='3' fill='var(--ink)' opacity='0.6'/><circle cx='10' cy='28' r='2.5' fill='var(--ink)' opacity='0.5'/><circle cx='88' cy='32' r='2.5' fill='var(--ink)' opacity='0.5'/></g></svg>`,

	"tree-hier": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<circle cx='50' cy='8' r='4' fill='var(--accent)'/>` +
		`<g stroke='var(--ink-3)' stroke-width='0.6'><line x1='50' y1='12' x2='22' y2='30'/><line x1='50' y1='12' x2='50' y2='30'/><line x1='50' y1='12' x2='78' y2='30'/></g>` +
		`<circle cx='22' cy='32' r='3' fill='var(--ink)' opacity='0.7'/><circle cx='50' cy='32' r='3' fill='var(--ink)' opacity='0.7'/><circle cx='78' cy='32' r='3' fill='var(--ink)' opacity='0.7'/>` +
		`<g stroke='var(--ink-3)' stroke-width='0.5' opacity='0.7'><line x1='22' y1='35' x2='14' y2='52'/><line x1='22' y1='35' x2='30' y2='52'/><line x1='50' y1='35' x2='42' y2='52'/><line x1='50' y1='35' x2='58' y2='52'/><line x1='78' y1='35' x2='70' y2='52'/><line x1='78' y1='35' x2='86' y2='52'/></g>` +
		`<g fill='var(--ink)' opacity='0.5'><circle cx='14' cy='54' r='2'/><circle cx='30' cy='54' r='2'/><circle cx='42' cy='54' r='2'/><circle cx='58' cy='54' r='2'/><circle cx='70' cy='54' r='2'/><circle cx='86' cy='54' r='2'/></g></svg>`,

	"radial-burst": `<svg viewBox='0 0 100 60' style='width:75%;height:auto'>` +
		`<circle cx='50' cy='30' r='22' fill='none' stroke='var(--rule-2)' stroke-width='0.5' stroke-dasharray='2 2'/>` +
		`<g stroke='var(--ink-3)' stroke-width='0.5'><line x1='50' y1='30' x2='72' y2='30'/><line x1='50' y1='30' x2='61' y2='49'/><line x1='50' y1='30' x2='39' y2='49'/><line x1='50' y1='30' x2='28' y2='30'/><line x1='50' y1='30' x2='39' y2='11'/><line x1='50' y1='30' x2='61' y2='11'/></g>` +
		`<g fill='var(--ink)' opacity='0.7'><circle cx='72' cy='30' r='2.5'/><circle cx='61' cy='49' r='2.5'/><circle cx='39' cy='49' r='2.5'/><circle cx='28' cy='30' r='2.5'/><circle cx='39' cy='11' r='2.5'/><circle cx='61' cy='11' r='2.5'/></g>` +
		`<circle cx='50' cy='30' r='5' fill='var(--accent)'/></svg>`,

	"matrix-heatmap": `<svg viewBox='0 0 100 60' style='width:84%;height:auto'><g>` +
		matrixGrid() + `</g></svg>`,

	// ── metaphors ────────────────────────────────────────────────────
	"plant-grow": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<line x1='4' y1='52' x2='96' y2='52' stroke='var(--ink)' stroke-width='1'/>` +
		`<line x1='20' y1='52' x2='20' y2='16' stroke='var(--ink)' stroke-width='1.2'/>` +
		`<path d='M 20 30 Q 11 26 8 18' fill='none' stroke='var(--accent)' stroke-width='1.4'/>` +
		`<path d='M 20 24 Q 30 20 33 12' fill='none' stroke='var(--accent)' stroke-width='1.4'/>` +
		`<circle cx='20' cy='14' r='4' fill='var(--accent)'/>` +
		`<line x1='50' y1='52' x2='50' y2='30' stroke='var(--ink)' stroke-width='1.2'/>` +
		`<path d='M 50 38 Q 42 34 40 28' fill='none' stroke='var(--accent)' stroke-width='1.2'/>` +
		`<circle cx='50' cy='28' r='3' fill='var(--accent)'/>` +
		`<line x1='80' y1='52' x2='80' y2='44' stroke='var(--ink)' stroke-width='1'/>` +
		`<path d='M 80 46 Q 76 43 74 41' fill='none' stroke='var(--accent)' stroke-width='1'/>` +
		`<path d='M 80 46 Q 84 43 86 41' fill='none' stroke='var(--accent)' stroke-width='1'/></svg>`,

	"garden": `<svg viewBox='0 0 100 60' style='width:84%;height:auto'>` +
		`<line x1='2' y1='52' x2='98' y2='52' stroke='var(--ink)' stroke-width='1'/>` +
		`<g fill='var(--accent)'><circle cx='10' cy='44' r='5'/><circle cx='22' cy='40' r='6'/><circle cx='35' cy='42' r='5'/></g>` +
		`<g fill='var(--ink)' opacity='0.6'><circle cx='50' cy='46' r='3'/><circle cx='58' cy='44' r='4'/></g>` +
		`<g fill='var(--accent)' opacity='0.55'><circle cx='75' cy='42' r='5'/><circle cx='86' cy='44' r='4'/><circle cx='94' cy='46' r='3'/></g>` +
		`<g stroke='var(--ink)' stroke-width='0.7'><line x1='10' y1='49' x2='10' y2='52'/><line x1='22' y1='46' x2='22' y2='52'/><line x1='35' y1='47' x2='35' y2='52'/><line x1='50' y1='49' x2='50' y2='52'/><line x1='58' y1='48' x2='58' y2='52'/><line x1='75' y1='47' x2='75' y2='52'/><line x1='86' y1='48' x2='86' y2='52'/></g></svg>`,

	"paper-fold": `<svg viewBox='0 0 100 60' style='width:70%;height:auto'>` +
		`<path d='M 20 10 L 70 10 L 80 20 L 80 50 L 20 50 Z' fill='var(--surface)' stroke='var(--ink)' stroke-width='1.2'/>` +
		`<path d='M 70 10 L 70 20 L 80 20' fill='none' stroke='var(--ink)' stroke-width='1.2'/>` +
		`<line x1='28' y1='24' x2='66' y2='24' stroke='var(--ink-3)' stroke-width='0.6'/>` +
		`<line x1='28' y1='30' x2='72' y2='30' stroke='var(--ink-3)' stroke-width='0.6'/>` +
		`<line x1='28' y1='36' x2='60' y2='36' stroke='var(--ink-3)' stroke-width='0.6'/>` +
		`<line x1='28' y1='42' x2='68' y2='42' stroke='var(--ink-3)' stroke-width='0.6'/></svg>`,

	// ── misc ─────────────────────────────────────────────────────────
	"cmd-palette": `<svg viewBox='0 0 100 60' style='width:75%;height:auto'>` +
		`<rect x='32' y='20' width='14' height='18' rx='2' fill='var(--surface)' stroke='var(--ink)' stroke-width='1.4'/>` +
		`<text x='39' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='10' font-weight='700' fill='var(--ink)'>⌘</text>` +
		`<rect x='54' y='20' width='14' height='18' rx='2' fill='var(--surface)' stroke='var(--ink)' stroke-width='1.4'/>` +
		`<text x='61' y='33' text-anchor='middle' font-family='IBM Plex Mono' font-size='10' font-weight='700' fill='var(--ink)'>K</text>` +
		`<text x='50' y='52' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' fill='var(--ink-3)' letter-spacing='0.08em'>FIND ANYTHING</text></svg>`,

	"text-cursor": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<rect x='10' y='24' width='80' height='14' fill='var(--surface)' stroke='var(--ink)' stroke-width='1'/>` +
		`<text x='14' y='34' font-family='IBM Plex Mono' font-size='8' fill='var(--ink)'>type here</text>` +
		`<line x1='52' y1='27' x2='52' y2='35' stroke='var(--accent)' stroke-width='1.4'/></svg>`,

	"file-icons": `<svg viewBox='0 0 100 60' style='width:80%;height:auto'>` +
		`<g font-family='IBM Plex Mono' font-size='8' fill='var(--ink-2)'>` +
		`<g transform='translate(10,12)'><path d='M 0 0 L 8 0 L 11 3 L 11 14 L 0 14 Z' fill='var(--surface)' stroke='var(--accent)' stroke-width='1'/><path d='M 8 0 L 8 3 L 11 3' fill='none' stroke='var(--accent)' stroke-width='1'/><text x='14' y='10'>spec.md</text></g>` +
		`<g transform='translate(10,28)'><path d='M 0 0 L 8 0 L 11 3 L 11 14 L 0 14 Z' fill='var(--surface)' stroke='var(--accent)' stroke-width='1'/><path d='M 8 0 L 8 3 L 11 3' fill='none' stroke='var(--accent)' stroke-width='1'/><text x='14' y='10'>ideas.md</text></g>` +
		`<g transform='translate(10,44)'><path d='M 0 0 L 8 0 L 11 3 L 11 14 L 0 14 Z' fill='var(--surface)' stroke='var(--accent)' stroke-width='1'/><path d='M 8 0 L 8 3 L 11 3' fill='none' stroke='var(--accent)' stroke-width='1'/><text x='14' y='10'>todo.md</text></g>` +
		`</g></svg>`,
}

// ─── helper builders (kept here to keep the registry compact) ───

// gridCells builds N×M rect cells, packed in the upper area, with some
// indices marked "filled" (var(--accent)) — used by calendar-grid.
func gridCells(cols, rows int, startX, startY int, cellW, cellH int, filled []int) string {
	fm := map[int]bool{}
	for _, i := range filled {
		fm[i] = true
	}
	var sb stringBuilder
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := r*cols + c
			x := startX + c*(cellW+1)
			y := startY + r*(cellH+1)
			if fm[idx] {
				sb.Appendf(`<rect x='%d' y='%d' width='%d' height='%d' fill='var(--accent)'/>`, x, y, cellW, cellH)
			} else {
				sb.Appendf(`<rect x='%d' y='%d' width='%d' height='%d' fill='none' stroke='var(--rule-2)' stroke-width='0.4'/>`, x, y, cellW, cellH)
			}
		}
	}
	return sb.String()
}

// yearGrid is a deterministic-noise 12×20 grid for year-heatmap.
func yearGrid() string {
	var sb stringBuilder
	for r := 0; r < 7; r++ {
		for c := 0; c < 26; c++ {
			x := 4.0 + float64(c)*3.5
			y := 6.0 + float64(r)*7
			seed := (r*113 + c*37 + 41) % 100
			switch {
			case seed < 32:
				sb.Appendf(`<rect x='%.2f' y='%.2f' width='3' height='6' fill='var(--accent)'/>`, x, y)
			case seed < 55:
				sb.Appendf(`<rect x='%.2f' y='%.2f' width='3' height='6' fill='var(--accent)' opacity='0.45'/>`, x, y)
			default:
				sb.Appendf(`<rect x='%.2f' y='%.2f' width='3' height='6' fill='none' stroke='var(--rule-2)' stroke-width='0.3'/>`, x, y)
			}
		}
	}
	return sb.String()
}

// matrixGrid is a deterministic 12×7 heatmap for matrix-heatmap.
func matrixGrid() string {
	var sb stringBuilder
	for r := 0; r < 7; r++ {
		for c := 0; c < 12; c++ {
			x := 8 + c*7
			y := 6 + r*7
			seed := (r*97 + c*23 + 17) % 100
			switch {
			case seed < 30:
				sb.Appendf(`<rect x='%d' y='%d' width='6' height='6' fill='var(--accent)'/>`, x, y)
			case seed < 55:
				sb.Appendf(`<rect x='%d' y='%d' width='6' height='6' fill='var(--accent)' opacity='0.45'/>`, x, y)
			default:
				sb.Appendf(`<rect x='%d' y='%d' width='6' height='6' fill='none' stroke='var(--rule-2)' stroke-width='0.3'/>`, x, y)
			}
		}
	}
	return sb.String()
}

// phoneScreen wraps a body fragment in a phone-frame at viewBox 0 0 100 60.
// Body is positioned inside a 36×42 phone screen centered horizontally.
func phoneScreen(body string) string {
	return `<svg viewBox='0 0 100 60' style='width:55%;height:auto'><g transform='translate(30,4)'>` +
		`<rect x='0' y='0' width='42' height='52' rx='4' fill='var(--surface)' stroke='var(--ink)' stroke-width='1.2'/>` +
		`<line x1='15' y1='4' x2='27' y2='4' stroke='var(--ink-3)' stroke-width='0.5'/>` +
		body +
		`</g></svg>`
}
