//go:generate go run ./cmd/gen-shortcodes -in emojis.json -out shortcodes_gen.go

package demoji

// builtinEmoticons is the default emoticon → emoji table.
// Order matters: the first active entry for a given emoji becomes the canonical
// text representation when converting unicode → emoticons (demoji mode).
// Longer variants are listed before shorter ones so sort-by-length produces
// the expected canonical ranking automatically.
var builtinEmoticons = [][2]string{
	// Angel / innocent
	{"O:-)", "😇"},
	{"O:)", "😇"},

	// Angry
	{">:-(", "😠"},
	{">:(", "😠"},

	// Crying
	{":'(", "😢"},

	// Confused
	{":-/", "😕"},
	{":-\\", "😕"},
	{":/", "😕"},

	// Cool / sunglasses
	{"B-)", "😎"},
	{"8-)", "😎"},
	{"B)", "😎"},
	{"8)", "😎"},

	// Happy / smile
	{":-)", "🙂"},
	{":)", "🙂"},

	// Big grin
	{":-D", "😀"},
	{":D", "😀"},

	// Sad
	{":-(", "🙁"},
	{":(", "🙁"},

	// Wink
	{";-)", "😉"},
	{";)", "😉"},

	// Tongue + wink
	{";-P", "😜"},
	{";P", "😜"},

	// Tongue
	{":-P", "😛"},
	{":-p", "😛"},
	{":P", "😛"},
	{":p", "😛"},

	// Surprised / open mouth
	{":-O", "😮"},
	{":-o", "😮"},
	{":O", "😮"},
	{":o", "😮"},

	// Neutral / meh
	{":-|", "😐"},
	{":|", "😐"},

	// Kiss
	{":-*", "😘"},
	{":*", "😘"},

	// Expressionless
	{"-_-", "😑"},

	// Laughing / XD
	{"XD", "😆"},
	{"xD", "😆"},

	// Love
	{"<3", "❤️"},

	// Broken heart
	{"</3", "💔"},
}
