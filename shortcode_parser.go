package demoji

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// shortcodeInlineParser is a goldmark inline parser that converts :shortcode:
// patterns to unicode emoji. It must be an inline parser (not an AST
// transformer) because shortcode names can contain underscores, and goldmark's
// emphasis tokenizer fragments text at '_' boundaries before any AST
// transformer runs — meaning a transformer would never see the full string
// ":slightly_smiling_face:" in a single text node.
//
// By registering ':' as a trigger character, we consume the complete
// :name: token before the emphasis parser gets a chance to look inside it.
type shortcodeInlineParser struct {
	table map[string]string // :shortcode: → emoji
}

func (p *shortcodeInlineParser) Trigger() []byte {
	return []byte{':'}
}

func (p *shortcodeInlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	// Need at least ":x:" (3 bytes) and must start with the trigger.
	if len(line) < 3 || line[0] != ':' {
		return nil
	}
	// Scan the name: one or more valid shortcode characters.
	i := 1
	for i < len(line) && isShortcodeChar(line[i]) {
		i++
	}
	// Must have at least one name char and end with a closing ':'.
	if i < 2 || i >= len(line) || line[i] != ':' {
		return nil
	}
	shortcode := string(line[:i+1])
	emoji, ok := p.table[shortcode]
	if !ok {
		return nil
	}
	block.Advance(i + 1)
	return ast.NewString([]byte(emoji))
}

// isShortcodeChar reports whether b is valid inside a GitHub shortcode name.
func isShortcodeChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') ||
		b == '_' || b == '-' || b == '+'
}
