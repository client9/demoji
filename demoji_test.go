package demoji_test

import (
	"maps"
	"testing"

	"github.com/client9/demoji"
)

// --- unicode → emoticon (default) ---

func TestUnicodeToEmoticon(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{"smile", "Hello 🙂 world", "Hello :-) world"},
		{"grin", "Woohoo 😀", "Woohoo :-D"},
		{"multiple", "😀 and 😉", ":-D and ;-)"},
		{"love", "I ❤️ Go", "I <3 Go"},
		{"broken heart", "💔", "</3"},
		{"no emoji", "Plain text", "Plain text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := demoji.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestVS16(t *testing.T) {
	// ❤️ is U+2764 U+FE0F — emoticons table stores the VS16 form, so it converts.
	got := demoji.Replace("I ❤️ Go")
	if got != "I <3 Go" {
		t.Errorf("VS16 emoticon: got %q", got)
	}
	// In unicode→shortcode direction, GitHub data stores the bare codepoint (U+2764).
	// The VS16 logic adds a U+FE0F variant, so ❤️ also maps to :heart:.
	r := demoji.New(demoji.Config{From: demoji.FormatUnicode, To: demoji.FormatShortcode})
	got = r.Replace("I ❤️ Go")
	if got != "I :heart: Go" {
		t.Errorf("VS16 shortcode: got %q", got)
	}
	got = r.Replace("I ❤ Go")
	if got != "I :heart: Go" {
		t.Errorf("bare shortcode: got %q", got)
	}
}

// --- unicode → shortcode ---

func TestUnicodeToShortcode(t *testing.T) {
	r := demoji.New(demoji.Config{
		From: demoji.FormatUnicode,
		To:   demoji.FormatShortcode,
	})
	tests := []struct{ name, in, want string }{
		{"smile", "Hello 🙂 world", "Hello :slightly_smiling_face: world"},
		{"grin", "Woohoo 😀", "Woohoo :grinning:"},
		{"heart", "I ❤️ Go", "I :heart: Go"},
		{"no emoji", "Plain text", "Plain text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- emoticon → unicode ---

func TestEmoticonToUnicode(t *testing.T) {
	r := demoji.New(demoji.Config{From: demoji.FormatEmoticon})
	tests := []struct{ name, in, want string }{
		{"classic smile", "Hello :-) world", "Hello 🙂 world"},
		{"short smile", "Hi :)", "Hi 🙂"},
		{"wink", ";-) wink", "😉 wink"},
		{"longer first", "O:-) angel", "😇 angel"},
		{"love", "I <3 Go", "I ❤️ Go"},
		{"multiple", ":-) and :-(", "🙂 and 🙁"},
		{"no emoticons", "Plain text", "Plain text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- shortcode → unicode ---

func TestShortcodeToUnicode(t *testing.T) {
	r := demoji.New(demoji.Config{From: demoji.FormatShortcode})
	tests := []struct{ name, in, want string }{
		{"smile", "Hello :slightly_smiling_face: world", "Hello 🙂 world"},
		{"grin", ":grinning:", "😀"},
		{"heart", "I :heart: Go", "I ❤ Go"},
		{"no shortcodes", "Plain text", "Plain text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- emoticon | shortcode → unicode (single pass, bitfield) ---

func TestCombinedToUnicode(t *testing.T) {
	r := demoji.New(demoji.Config{From: demoji.FormatEmoticon | demoji.FormatShortcode})
	tests := []struct{ name, in, want string }{
		{"both", ":-) and :grinning: are both happy", "🙂 and 😀 are both happy"},
		{"shortcode only", ":wink:", "😉"},
		{"emoticon only", ";-)", "😉"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- zero Config is a no-op ---

func TestNoop(t *testing.T) {
	r := demoji.New(demoji.Config{})
	got := r.Replace("Hello 🙂 :-) :grinning:")
	if got != "Hello 🙂 :-) :grinning:" {
		t.Errorf("Config{} should be no-op, got %q", got)
	}
}

// --- ReplaceBytes ---

func TestReplaceBytes(t *testing.T) {
	got := demoji.ReplaceBytes([]byte("Hello 🙂"))
	if string(got) != "Hello :-)" {
		t.Errorf("got %q", got)
	}
}

// --- custom maps ---

func TestCustomEmoticons(t *testing.T) {
	custom := maps.Clone(demoji.DefaultEmoticons())
	custom["^^"] = "😊"
	r := demoji.New(demoji.Config{
		From:      demoji.FormatEmoticon,
		Emoticons: custom,
	})
	got := r.Replace("Hello ^^ world")
	if got != "Hello 😊 world" {
		t.Errorf("got %q", got)
	}
}

func TestExcludeEmoticon(t *testing.T) {
	// Remove both forms of the smile emoticon — emoji should pass through.
	custom := maps.Clone(demoji.DefaultEmoticons())
	delete(custom, ":-)") // remove canonical
	delete(custom, ":)")  // remove fallback
	r := demoji.New(demoji.Config{
		From:      demoji.FormatUnicode,
		To:        demoji.FormatEmoticon,
		Emoticons: custom,
	})
	got := r.Replace("Hello 🙂 world")
	if got != "Hello 🙂 world" {
		t.Errorf("got %q", got)
	}
}

func TestCanonicalFallback(t *testing.T) {
	// Remove the canonical ":-)" — should fall back to ":)".
	custom := maps.Clone(demoji.DefaultEmoticons())
	delete(custom, ":-)")
	r := demoji.New(demoji.Config{
		From:      demoji.FormatUnicode,
		To:        demoji.FormatEmoticon,
		Emoticons: custom,
	})
	got := r.Replace("Hello 🙂 world")
	if got != "Hello :) world" {
		t.Errorf("got %q", got)
	}
}

func TestExcludeShortcode(t *testing.T) {
	custom := maps.Clone(demoji.DefaultShortcodes())
	delete(custom, ":grinning:")
	r := demoji.New(demoji.Config{
		From:       demoji.FormatUnicode,
		To:         demoji.FormatShortcode,
		Shortcodes: custom,
	})
	// 😀 excluded → passes through unchanged
	got := r.Replace("Woohoo 😀")
	if got != "Woohoo 😀" {
		t.Errorf("got %q", got)
	}
}

func BenchmarkNew(b *testing.B) {
	b.Run("default", func(b *testing.B) {
		for range b.N {
			_ = demoji.New(demoji.DefaultConfig())
		}
	})
	b.Run("with_custom_map", func(b *testing.B) {
		extra := maps.Clone(demoji.DefaultEmoticons())
		extra["^^"] = "😊"
		cfg := demoji.Config{From: demoji.FormatEmoticon, Emoticons: extra}
		for range b.N {
			_ = demoji.New(cfg)
		}
	})
}

func BenchmarkReplace(b *testing.B) {
	r := demoji.New(demoji.DefaultConfig())
	s := "Hello 🙂 and 😀 and ❤️ Go"
	b.ResetTimer()
	for range b.N {
		_ = r.Replace(s)
	}
}
