package demoji

import (
	"sort"
	"strings"
)

// vs16 is U+FE0F (variation selector-16), the emoji-presentation variant selector.
const vs16 = "️"

// buildPairs assembles the (source → target) replacement table from cfg.
// Pairs are sorted longest-source-first so that longer patterns (e.g. "O:-)")
// are matched before shorter prefixes (":-)").
func buildPairs(cfg Config) [][2]string {
	var pairs [][2]string

	if cfg.From&FormatUnicode != 0 {
		// unicode → target text
		switch cfg.To {
		case FormatEmoticon:
			pairs = unicodeToTextPairs(cfg.Emoticons, builtinEmoticons)
		case FormatShortcode:
			pairs = unicodeToTextPairs(cfg.Shortcodes, builtinShortcodes)
		}
	} else {
		// text → unicode; From is a bitfield
		if cfg.From&FormatEmoticon != 0 {
			for emoticon, emoji := range cfg.Emoticons {
				pairs = append(pairs, [2]string{emoticon, emoji})
			}
		}
		if cfg.From&FormatShortcode != 0 {
			for shortcode, emoji := range cfg.Shortcodes {
				pairs = append(pairs, [2]string{shortcode, emoji})
			}
		}
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

// unicodeToTextPairs builds (emoji → text) pairs for unicode→text directions.
// The canonical text for each emoji is the first active entry in builtins;
// user-added entries not in builtins claim any unclaimed emoji.
// VS16 variants are added for each mapped emoji not already ending in U+FE0F,
// ensuring ❤️ (U+2764 U+FE0F) is caught before the bare ❤ replacement fires.
func unicodeToTextPairs(active map[string]string, builtins [][2]string) [][2]string {
	canonical := make(map[string]string)

	for _, pair := range builtins {
		text := pair[0]
		currentEmoji, ok := active[text]
		if !ok {
			continue // excluded by caller
		}
		emoji := currentEmoji // may be overridden
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

	// Add VS16 variants.
	for emoji, text := range canonical {
		if !strings.HasSuffix(emoji, vs16) {
			if variant := emoji + vs16; canonical[variant] == "" {
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
