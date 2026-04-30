// Package demoji provides a goldmark extension that converts between unicode
// emoji, classic ASCII emoticons, and GitHub-style shortcodes.
//
// The source and target formats are specified independently. Multiple source
// formats can be OR'd together for a single-pass conversion:
//
//	// unicode emoji → ASCII emoticons (default)
//	demoji.New()
//
//	// ASCII emoticons → unicode emoji
//	demoji.New(demoji.WithFrom(demoji.FormatEmoticon))
//
//	// both emoticons and shortcodes → unicode, in one pass
//	demoji.New(demoji.WithFrom(demoji.FormatEmoticon | demoji.FormatShortcode))
//
//	// unicode → GitHub shortcodes
//	demoji.New(demoji.WithFrom(demoji.FormatUnicode), demoji.WithTo(demoji.FormatShortcode))
package demoji

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// Format identifies an emoji representation. It is a bitfield so multiple
// source formats can be combined with WithFrom.
type Format uint

const (
	FormatUnicode   Format = 1 << iota // 😀  actual unicode emoji codepoints
	FormatEmoticon                      // :-)  classic ASCII emoticons
	FormatShortcode                     // :grinning:  GitHub / Gemoji shortcodes
)

// Option configures the extension.
type Option func(*config)

type config struct {
	from       Format
	to         Format
	emoticons  map[string]string // emoticon → emoji
	shortcodes map[string]string // :shortcode: → emoji
}

// WithFrom sets the source format(s). OR multiple values to convert several
// formats in a single AST pass:
//
//	demoji.WithFrom(demoji.FormatEmoticon | demoji.FormatShortcode)
//
// Default: FormatUnicode.
func WithFrom(f Format) Option {
	return func(c *config) { c.from = f }
}

// WithTo sets the target format. Must be a single Format value.
// Default: FormatEmoticon.
func WithTo(f Format) Option {
	return func(c *config) { c.to = f }
}

// WithAdditional adds or overrides mappings for the given format.
// Keys are source patterns; values are the corresponding unicode emoji strings.
//
//	// add a custom emoticon
//	demoji.WithAdditional(demoji.FormatEmoticon, map[string]string{"^^": "😊"})
//
//	// add a custom shortcode
//	demoji.WithAdditional(demoji.FormatShortcode, map[string]string{":custom:": "🦄"})
func WithAdditional(format Format, m map[string]string) Option {
	return func(c *config) {
		switch format {
		case FormatEmoticon:
			for k, v := range m {
				c.emoticons[k] = v
			}
		case FormatShortcode:
			for k, v := range m {
				c.shortcodes[k] = v
			}
		}
	}
}

// WithExclude removes the named patterns from the active mapping for the given format.
//
//	demoji.WithExclude(demoji.FormatEmoticon, ":-)", ":)")
//	demoji.WithExclude(demoji.FormatShortcode, ":cry:")
func WithExclude(format Format, keys ...string) Option {
	return func(c *config) {
		switch format {
		case FormatEmoticon:
			for _, k := range keys {
				delete(c.emoticons, k)
			}
		case FormatShortcode:
			for _, k := range keys {
				delete(c.shortcodes, k)
			}
		}
	}
}

// Extension is the goldmark extension for emoji format conversion.
type Extension struct {
	cfg *config
}

// New creates the extension. With no options the default behavior is
// unicode emoji → ASCII emoticons (equivalent to the old ModeDemoji).
func New(opts ...Option) *Extension {
	cfg := &config{
		from:       FormatUnicode,
		to:         FormatEmoticon,
		emoticons:  make(map[string]string, len(builtinEmoticons)),
		shortcodes: make(map[string]string, len(builtinShortcodes)),
	}
	for _, pair := range builtinEmoticons {
		cfg.emoticons[pair[0]] = pair[1]
	}
	for _, pair := range builtinShortcodes {
		cfg.shortcodes[pair[0]] = pair[1]
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return &Extension{cfg: cfg}
}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(m goldmark.Markdown) {
	// When converting shortcodes → unicode, use an inline parser rather than
	// an AST transformer. Shortcode names can contain underscores, and
	// goldmark's emphasis tokenizer fragments text at '_' boundaries before
	// any AST transformer runs, so a transformer would never see the full
	// ":slightly_smiling_face:" string in a single node. The inline parser
	// triggers on ':' and consumes the whole :name: before that happens.
	if e.cfg.from&FormatShortcode != 0 && e.cfg.from&FormatUnicode == 0 {
		m.Parser().AddOptions(
			parser.WithInlineParsers(
				util.Prioritized(&shortcodeInlineParser{table: e.cfg.shortcodes}, 999),
			),
		)
	}
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&emojiTransformer{cfg: e.cfg}, 100),
		),
	)
}
