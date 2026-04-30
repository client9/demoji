package demoji_test

import (
	"bytes"
	"testing"

	"github.com/client9/demoji"
	"github.com/yuin/goldmark"
)

func render(t *testing.T, ext *demoji.Extension, src string) string {
	t.Helper()
	md := goldmark.New(goldmark.WithExtensions(ext))
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		t.Fatalf("Convert: %v", err)
	}
	return buf.String()
}

// --- unicode → emoticon (default) ---

func TestDemojiToEmoticon(t *testing.T) {
	ext := demoji.New() // FormatUnicode → FormatEmoticon by default

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"smile", "Hello 🙂 world", "<p>Hello :-) world</p>\n"},
		{"grin", "Woohoo 😀", "<p>Woohoo :-D</p>\n"},
		{"multiple", "😀 and 😉", "<p>:-D and ;-)</p>\n"},
		{"love", "I ❤️ Go", "<p>I &lt;3 Go</p>\n"},
		{"broken heart", "💔", "<p>&lt;/3</p>\n"},
		{"no emoji", "Plain text", "<p>Plain text</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ext, tt.input)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- unicode → shortcode ---

func TestDemojiToShortcode(t *testing.T) {
	ext := demoji.New(
		demoji.WithFrom(demoji.FormatUnicode),
		demoji.WithTo(demoji.FormatShortcode),
	)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"smile", "Hello 🙂 world", "<p>Hello :slightly_smiling_face: world</p>\n"},
		{"grin", "Woohoo 😀", "<p>Woohoo :grinning:</p>\n"},
		{"wink", "😉", "<p>:wink:</p>\n"},
		{"heart", "I ❤️ Go", "<p>I :heart: Go</p>\n"},
		{"no emoji", "Plain text", "<p>Plain text</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ext, tt.input)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- emoticon → unicode ---

func TestRemojiFromEmoticon(t *testing.T) {
	ext := demoji.New(demoji.WithFrom(demoji.FormatEmoticon))

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"classic smile", "Hello :-) world", "<p>Hello 🙂 world</p>\n"},
		{"short smile", "Hi :)", "<p>Hi 🙂</p>\n"},
		{"wink", ";-) wink", "<p>😉 wink</p>\n"},
		{"angel longer first", "O:-) angel", "<p>😇 angel</p>\n"},
		{"love", "I <3 Go", "<p>I ❤️ Go</p>\n"},
		{"multiple", ":-) and :-(", "<p>🙂 and 🙁</p>\n"},
		{"no emoticons", "Plain text", "<p>Plain text</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ext, tt.input)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- shortcode → unicode ---

func TestRemojiFromShortcode(t *testing.T) {
	ext := demoji.New(demoji.WithFrom(demoji.FormatShortcode))

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"smile", "Hello :slightly_smiling_face: world", "<p>Hello 🙂 world</p>\n"},
		{"grin", ":grinning:", "<p>😀</p>\n"},
		{"long shortcode first", ":stuck_out_tongue_winking_eye:", "<p>😜</p>\n"},
		{"heart", "I :heart: Go", "<p>I ❤ Go</p>\n"},
		{"no shortcodes", "Plain text", "<p>Plain text</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ext, tt.input)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- emoticon | shortcode → unicode (single pass, bitfield) ---

func TestRemojiCombined(t *testing.T) {
	ext := demoji.New(
		demoji.WithFrom(demoji.FormatEmoticon | demoji.FormatShortcode),
	)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			"emoticon and shortcode in same node",
			":-) and :grinning: are both happy",
			"<p>🙂 and 😀 are both happy</p>\n",
		},
		{
			"shortcode only",
			":wink:",
			"<p>😉</p>\n",
		},
		{
			"emoticon only",
			";-)",
			"<p>😉</p>\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ext, tt.input)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- code blocks and spans are skipped ---

func TestSkipsCodeSpan(t *testing.T) {
	tests := []struct {
		mode  demoji.Format
		input string
		want  string
	}{
		{
			demoji.FormatUnicode,
			"outside 🙂 but `inside 🙂 code`",
			"<p>outside :-) but <code>inside 🙂 code</code></p>\n",
		},
		{
			demoji.FormatEmoticon,
			"outside :-) but `inside :-)  code`",
			"<p>outside 🙂 but <code>inside :-)  code</code></p>\n",
		},
		{
			demoji.FormatShortcode,
			"outside :grinning: but `inside :grinning: code`",
			"<p>outside 😀 but <code>inside :grinning: code</code></p>\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input[:15], func(t *testing.T) {
			got := render(t, demoji.New(demoji.WithFrom(tt.mode)), tt.input)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestSkipsFencedCodeBlock(t *testing.T) {
	input := "```\n🙂 :-) :grinning:\n```"
	got := render(t, demoji.New(), input)
	want := "<pre><code>🙂 :-) :grinning:\n</code></pre>\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

// --- configuration: WithAdditional and WithExclude ---

func TestWithAdditional(t *testing.T) {
	ext := demoji.New(
		demoji.WithFrom(demoji.FormatEmoticon),
		demoji.WithAdditional(demoji.FormatEmoticon, map[string]string{"^^": "😊"}),
	)
	got := render(t, ext, "Hello ^^ world")
	want := "<p>Hello 😊 world</p>\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestWithAdditionalShortcode(t *testing.T) {
	ext := demoji.New(
		demoji.WithFrom(demoji.FormatShortcode),
		demoji.WithAdditional(demoji.FormatShortcode, map[string]string{":unicorn:": "🦄"}),
	)
	got := render(t, ext, "Hello :unicorn: world")
	want := "<p>Hello 🦄 world</p>\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestWithExclude(t *testing.T) {
	// Exclude both variants so 🙂 is no longer converted.
	ext := demoji.New(
		demoji.WithExclude(demoji.FormatEmoticon, ":-)", ":)"),
	)
	got := render(t, ext, "Hello 🙂 world")
	want := "<p>Hello 🙂 world</p>\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestWithExcludeCanonicalFallback(t *testing.T) {
	// Exclude canonical ":-)" but keep ":)"; demoji should fall back to ":)".
	ext := demoji.New(
		demoji.WithExclude(demoji.FormatEmoticon, ":-)"),
	)
	got := render(t, ext, "Hello 🙂 world")
	want := "<p>Hello :) world</p>\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestWithExcludeShortcode(t *testing.T) {
	ext := demoji.New(
		demoji.WithFrom(demoji.FormatUnicode),
		demoji.WithTo(demoji.FormatShortcode),
		demoji.WithExclude(demoji.FormatShortcode, ":grinning:"),
	)
	// 😀 excluded → passes through unchanged
	got := render(t, ext, "Woohoo 😀")
	want := "<p>Woohoo 😀</p>\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}
