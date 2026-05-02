// Package demoji converts between Unicode emoji, ASCII emoticons, and
// GitHub-style shortcodes in plain strings.
//
// Use package-level functions for the default conversion (unicode → emoticons):
//
//	clean := demoji.Replace(s)
//	clean := demoji.ReplaceBytes(b)
//
// Configure with a Config struct:
//
//	r := demoji.New(demoji.Config{
//	    From: demoji.FormatEmoticon,
//	})
//	clean := r.Replace(s)
//
// Zero From/To are not auto-resolved — use DefaultConfig() as a base:
//
//	cfg := demoji.DefaultConfig()
//	cfg.From = demoji.FormatShortcode
//	r := demoji.New(cfg)
package demoji

import (
	"strings"
	"sync"
)

// Format identifies an emoji representation. It is a bitfield so multiple
// source formats can be combined in the From field.
type Format uint

const (
	FormatUnicode   Format = 1 << iota // 😀  actual unicode emoji codepoints
	FormatEmoticon                     // :-)  classic ASCII emoticons
	FormatShortcode                    // :grinning:  GitHub / Gemoji shortcodes
)

// Config configures a Replacer.
type Config struct {
	From       Format            // source format(s); bitfield for multiple
	To         Format            // target format; only used when From includes FormatUnicode
	Emoticons  map[string]string // emoticon → unicode emoji; nil = use builtins
	Shortcodes map[string]string // :shortcode: → unicode emoji; nil = use builtins
}

// DefaultConfig returns the default configuration: unicode emoji → ASCII emoticons.
func DefaultConfig() Config {
	return Config{
		From: FormatUnicode,
		To:   FormatEmoticon,
	}
}

// Replacer applies emoji conversions. Create with New.
type Replacer struct {
	sr *strings.Replacer
}

var defaultReplacer = sync.OnceValue(func() *Replacer {
	return New(DefaultConfig())
})

// New creates a Replacer from cfg. Nil maps are resolved to the built-in
// tables. Zero From/To are not defaulted — New(Config{}) produces a no-op
// Replacer. Use DefaultConfig() for sensible starting values.
func New(cfg Config) *Replacer {
	if cfg.Emoticons == nil {
		cfg.Emoticons = getDefaultEmoticons()
	}
	if cfg.Shortcodes == nil {
		cfg.Shortcodes = getDefaultShortcodes()
	}
	pairs := buildPairs(cfg)
	args := make([]string, 0, len(pairs)*2)
	for _, p := range pairs {
		args = append(args, p[0], p[1])
	}
	return &Replacer{sr: strings.NewReplacer(args...)}
}

// Replace returns s with all active conversions applied.
func (r *Replacer) Replace(s string) string { return r.sr.Replace(s) }

// ReplaceBytes returns a copy of b with all active conversions applied.
func (r *Replacer) ReplaceBytes(b []byte) []byte { return []byte(r.sr.Replace(string(b))) }

// Replace returns s with the default conversion applied (unicode emoji → emoticons).
func Replace(s string) string { return defaultReplacer().Replace(s) }

// ReplaceBytes returns a copy of b with the default conversion applied.
func ReplaceBytes(b []byte) []byte { return defaultReplacer().ReplaceBytes(b) }

// DefaultEmoticons returns the shared built-in emoticon map. Read-only —
// call maps.Clone if you need an editable copy.
func DefaultEmoticons() map[string]string { return getDefaultEmoticons() }

// DefaultShortcodes returns the shared built-in shortcode map. Read-only —
// call maps.Clone if you need an editable copy.
func DefaultShortcodes() map[string]string { return getDefaultShortcodes() }

var (
	defaultEmoticonMap  map[string]string
	defaultShortcodeMap map[string]string
	emoticonsOnce       sync.Once
	shortcodesOnce      sync.Once
)

func getDefaultEmoticons() map[string]string {
	emoticonsOnce.Do(func() {
		defaultEmoticonMap = make(map[string]string, len(builtinEmoticons))
		for _, p := range builtinEmoticons {
			defaultEmoticonMap[p[0]] = p[1]
		}
	})
	return defaultEmoticonMap
}

func getDefaultShortcodes() map[string]string {
	shortcodesOnce.Do(func() {
		defaultShortcodeMap = make(map[string]string, len(builtinShortcodes))
		for _, p := range builtinShortcodes {
			defaultShortcodeMap[p[0]] = p[1]
		}
	})
	return defaultShortcodeMap
}
