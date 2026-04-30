package demoji

import (
	"bytes"
	"sort"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// vs16 is the Unicode variation selector-16 (U+FE0F), which follows many
// base emoji codepoints in their "emoji presentation" form (e.g. ❤ vs ❤️).
const vs16 = "️"

type emojiTransformer struct {
	cfg *config
}

func (t *emojiTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()
	pairs := buildPairs(t.cfg)
	if len(pairs) == 0 {
		return
	}
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindCodeBlock, ast.KindFencedCodeBlock, ast.KindCodeSpan,
			ast.KindHTMLBlock, ast.KindRawHTML:
			return ast.WalkSkipChildren, nil
		case ast.KindText:
			processTextNode(n.(*ast.Text), source, pairs)
		}
		return ast.WalkContinue, nil
	})
}

func processTextNode(node *ast.Text, source []byte, pairs [][2]string) {
	content := node.Segment.Value(source)
	newContent := applyReplacements(content, pairs)
	if bytes.Equal(content, newContent) {
		return
	}
	newNode := ast.NewString(newContent)
	node.Parent().ReplaceChild(node.Parent(), node, newNode)
}

func applyReplacements(content []byte, pairs [][2]string) []byte {
	for _, pair := range pairs {
		content = bytes.ReplaceAll(content, []byte(pair[0]), []byte(pair[1]))
	}
	return content
}

// buildPairs assembles the (source → target) replacement table from cfg.
// When from contains FormatUnicode, this is a demoji direction (unicode → text).
// Otherwise, all active text formats are merged into a single list targeting unicode.
// All tables are sorted longest-source-first to prevent shorter patterns from
// consuming the prefix of a longer one (e.g. ":)" before ":-)").
func buildPairs(cfg *config) [][2]string {
	var pairs [][2]string

	if cfg.from&FormatUnicode != 0 {
		// unicode → target text format
		switch cfg.to {
		case FormatEmoticon:
			pairs = unicodeToTextPairs(cfg.emoticons, builtinEmoticons)
		case FormatShortcode:
			pairs = unicodeToTextPairs(cfg.shortcodes, builtinShortcodes)
		}
	} else {
		// text format(s) → unicode; from is a bitfield.
		if cfg.from&FormatEmoticon != 0 {
			for emoticon, emoji := range cfg.emoticons {
				pairs = append(pairs, [2]string{emoticon, emoji})
			}
		}
		// FormatShortcode is intentionally omitted here: shortcodes are
		// handled by shortcodeInlineParser, which runs before this transformer
		// and correctly handles underscore-containing names that the emphasis
		// tokenizer would otherwise fragment across multiple text nodes.
	}

	sort.Slice(pairs, func(i, j int) bool {
		li, lj := len(pairs[i][0]), len(pairs[j][0])
		if li != lj {
			return li > lj
		}
		return pairs[i][0] < pairs[j][0]
	})
	return pairs
}

// unicodeToTextPairs builds (emoji → text) pairs for demoji/deshortcode mode.
// The canonical text for each emoji is the first active entry in builtins;
// user-added entries fill gaps for emoji not covered by builtins.
func unicodeToTextPairs(active map[string]string, builtins [][2]string) [][2]string {
	canonical := make(map[string]string) // emoji → canonical text

	for _, pair := range builtins {
		text, emoji := pair[0], pair[1]
		currentEmoji, ok := active[text]
		if !ok {
			continue // excluded by user
		}
		emoji = currentEmoji // may be overridden by user
		if _, seen := canonical[emoji]; !seen {
			canonical[emoji] = text
		}
	}

	// User-added entries not present in builtins claim any unclaimed emoji.
	builtinSet := make(map[string]bool, len(builtins))
	for _, p := range builtins {
		builtinSet[p[0]] = true
	}
	for text, emoji := range active {
		if !builtinSet[text] {
			if _, seen := canonical[emoji]; !seen {
				canonical[emoji] = text
			}
		}
	}

	// For each mapped emoji that does not already end in U+FE0F, also add the
	// variation-selector-16 form. Emoji keyboards typically insert the VS form
	// (e.g. "❤️" = U+2764 U+FE0F) even when the data source only records the
	// base codepoint (U+2764). Without this, replacing U+2764 → ":heart:"
	// would leave the dangling U+FE0F in the output.
	for emoji, text := range canonical {
		if !strings.HasSuffix(emoji, vs16) {
			variant := emoji + vs16
			if _, seen := canonical[variant]; !seen {
				canonical[variant] = text
			}
		}
	}

	pairs := make([][2]string, 0, len(canonical))
	for emoji, text := range canonical {
		pairs = append(pairs, [2]string{emoji, text})
	}
	return pairs
}
