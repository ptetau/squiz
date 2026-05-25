package renderer

// WFNaturalBox is the measured natural footprint for each WFLibrary
// entry, in viewBox units (the canonical 100×60 frame). Numbers come
// from the outer `<svg ... style='width:NN%;height:auto'>` attribute on
// each entry in wf.go — NN maps directly to viewBox-X because the
// viewBox is 100 wide; height stays at the full 60 since wireframes
// fill the vertical extent.
//
// Composing agents use this to size a `<use href="wf:…"/>` reference
// without guessing — e.g. `width = NaturalBox.W * scale`.
//
// Keys MUST match WFLibrary exactly; TestWFNaturalBoxCoverage catches
// drift in both directions.
var WFNaturalBox = map[string]NaturalBox{
	// ── calendars / dates ────────────────────────────────────────────
	"calendar-grid":  {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
	"calendar-week":  {W: 82, H: 60, AspectRatio: 82.0 / 60.0},
	"streak-counter": {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
	"day-strip":      {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
	"year-heatmap":   {W: 84, H: 60, AspectRatio: 84.0 / 60.0},
	"time-of-day":    {W: 78, H: 60, AspectRatio: 78.0 / 60.0},
	"clock":          {W: 70, H: 60, AspectRatio: 70.0 / 60.0},

	// ── charts ───────────────────────────────────────────────────────
	"spark-rising": {W: 82, H: 60, AspectRatio: 82.0 / 60.0},
	"spark-flat":   {W: 82, H: 60, AspectRatio: 82.0 / 60.0},
	"spark-noisy":  {W: 82, H: 60, AspectRatio: 82.0 / 60.0},
	"bars-up":      {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
	"donut":        {W: 70, H: 60, AspectRatio: 70.0 / 60.0},
	"gauge":        {W: 78, H: 60, AspectRatio: 78.0 / 60.0},
	"dot-trend":    {W: 84, H: 60, AspectRatio: 84.0 / 60.0},

	// ── identities / avatars ─────────────────────────────────────────
	"avatar-single":  {W: 60, H: 60, AspectRatio: 60.0 / 60.0},
	"avatar-pair":    {W: 60, H: 60, AspectRatio: 60.0 / 60.0},
	"avatar-circle":  {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
	"avatar-feed":    {W: 84, H: 60, AspectRatio: 84.0 / 60.0},
	"avatar-private": {W: 60, H: 60, AspectRatio: 60.0 / 60.0},

	// ── phone screens (all via phoneScreen helper → width:55%) ───────
	"phone-blank":   {W: 55, H: 60, AspectRatio: 55.0 / 60.0},
	"phone-list":    {W: 55, H: 60, AspectRatio: 55.0 / 60.0},
	"phone-card":    {W: 55, H: 60, AspectRatio: 55.0 / 60.0},
	"phone-input":   {W: 55, H: 60, AspectRatio: 55.0 / 60.0},
	"phone-tabs":    {W: 55, H: 60, AspectRatio: 55.0 / 60.0},
	"phone-onboard": {W: 55, H: 60, AspectRatio: 55.0 / 60.0},
	"phone-stats":   {W: 55, H: 60, AspectRatio: 55.0 / 60.0},

	// ── controls ─────────────────────────────────────────────────────
	"toggle-on":     {W: 55, H: 60, AspectRatio: 55.0 / 60.0},
	"toggle-off":    {W: 55, H: 60, AspectRatio: 55.0 / 60.0},
	"button-accent": {W: 70, H: 60, AspectRatio: 70.0 / 60.0},
	"button-ghost":  {W: 70, H: 60, AspectRatio: 70.0 / 60.0},
	"slider":        {W: 78, H: 60, AspectRatio: 78.0 / 60.0},
	"dropdown":      {W: 78, H: 60, AspectRatio: 78.0 / 60.0},

	// ── status ───────────────────────────────────────────────────────
	"badge-new":   {W: 60, H: 60, AspectRatio: 60.0 / 60.0},
	"pill-row":    {W: 84, H: 60, AspectRatio: 84.0 / 60.0},
	"snowflake":   {W: 55, H: 60, AspectRatio: 55.0 / 60.0},
	"lock":        {W: 45, H: 60, AspectRatio: 45.0 / 60.0},
	"check-large": {W: 50, H: 60, AspectRatio: 50.0 / 60.0},

	// ── typography samples ───────────────────────────────────────────
	"serif-sample": {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
	"sans-sample":  {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
	"mono-sample":  {W: 80, H: 60, AspectRatio: 80.0 / 60.0},

	// ── connections / graphs ─────────────────────────────────────────
	"graph-force":    {W: 84, H: 60, AspectRatio: 84.0 / 60.0},
	"tree-hier":      {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
	"radial-burst":   {W: 75, H: 60, AspectRatio: 75.0 / 60.0},
	"matrix-heatmap": {W: 84, H: 60, AspectRatio: 84.0 / 60.0},

	// ── metaphors ────────────────────────────────────────────────────
	"plant-grow": {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
	"garden":     {W: 84, H: 60, AspectRatio: 84.0 / 60.0},
	"paper-fold": {W: 70, H: 60, AspectRatio: 70.0 / 60.0},

	// ── misc ─────────────────────────────────────────────────────────
	"cmd-palette": {W: 75, H: 60, AspectRatio: 75.0 / 60.0},
	"text-cursor": {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
	"file-icons":  {W: 80, H: 60, AspectRatio: 80.0 / 60.0},
}
