# demoji
[![Go Reference](https://pkg.go.dev/badge/github.com/client9/demoji.svg)](https://pkg.go.dev/github.com/client9/demoji)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/client9/demoji/actions/workflows/go.yml/badge.svg)](https://github.com/client9/demoji/actions)

Converts between Unicode emoji, ASCII emoticons, and GitHub-style shortcodes in plain
strings — in either direction.

```
😀  ↔  :-D                       (unicode ↔ emoticon)
🙂  ↔  :slightly_smiling_face:   (unicode ↔ shortcode)
```

Zero dependencies. For a goldmark extension see
[goldmark-demoji](https://github.com/client9/goldmark-demoji).

## Install

```bash
go get github.com/client9/demoji
```

## Quick start

```go
import "github.com/client9/demoji"

// Package-level default: unicode emoji → ASCII emoticons.
clean := demoji.Replace(s)
clean  = demoji.ReplaceBytes(b)

// Configured instance.
r := demoji.New(demoji.Config{From: demoji.FormatEmoticon})
clean  = r.Replace(s)
```

## Formats

| Constant | Example | Description |
|----------|---------|-------------|
| `FormatUnicode` | `😀` | Unicode emoji codepoints |
| `FormatEmoticon` | `:-D` | Classic ASCII emoticons |
| `FormatShortcode` | `:grinning:` | GitHub / Gemoji shortcodes |

## Conversion directions

| Config | Effect |
|--------|--------|
| `DefaultConfig()` | unicode → emoticon (default) |
| `Config{From: FormatUnicode, To: FormatShortcode}` | unicode → shortcode |
| `Config{From: FormatEmoticon}` | emoticon → unicode |
| `Config{From: FormatShortcode}` | shortcode → unicode |
| `Config{From: FormatEmoticon \| FormatShortcode}` | both → unicode in one pass |
| `Config{}` | no-op |

## Configuration

### Config struct

```go
type Config struct {
    From       Format            // source format(s); bitfield for multiple
    To         Format            // target; only used when From includes FormatUnicode
    Emoticons  map[string]string // nil = use builtins
    Shortcodes map[string]string // nil = use builtins
}
```

### DefaultConfig

`New(Config{})` is a no-op — zero `From` means no conversions. Use `DefaultConfig()`
as a starting point when you want the defaults plus modifications:

```go
cfg := demoji.DefaultConfig()          // unicode → emoticon
cfg.From = demoji.FormatShortcode      // change direction
r := demoji.New(cfg)
```

### Custom maps

`DefaultEmoticons()` and `DefaultShortcodes()` return the shared built-in maps
(read-only). Clone before modifying:

```go
import "maps"

// Add a custom emoticon.
emoticons := maps.Clone(demoji.DefaultEmoticons())
emoticons["^^"] = "😊"
r := demoji.New(demoji.Config{
    From:      demoji.FormatEmoticon,
    Emoticons: emoticons,
})

// Exclude an emoticon.
emoticons = maps.Clone(demoji.DefaultEmoticons())
delete(emoticons, ":-)")
delete(emoticons, ":)")
r = demoji.New(demoji.Config{
    From:      demoji.FormatUnicode,
    To:        demoji.FormatEmoticon,
    Emoticons: emoticons,
})
```

### Canonical text selection

In unicode→text mode each emoji maps to one canonical text form — the first active entry
in the built-in table. To shift the canonical, exclude the current leader:

```go
// default: 🙂 → :-)
// after removing ":-)", the next entry ":)" becomes canonical:
emoticons := maps.Clone(demoji.DefaultEmoticons())
delete(emoticons, ":-)")
r := demoji.New(demoji.Config{
    From:      demoji.FormatUnicode,
    To:        demoji.FormatEmoticon,
    Emoticons: emoticons,
})
r.Replace("🙂") // → ":)"
```

## Built-in data

**Emoticons** (~40 entries): `O:-)` `>:(` `:'(` `:-)` `:)` `:-D` `:D` `:-(` `:(` `;-)`
`;)` `:-P` `:-O` `:-|` `:-*` `-_-` `XD` `<3` `</3` and variants.

**Shortcodes** (1913 entries): the full [GitHub emoji](https://api.github.com/emojis) set,
minus 23 GitHub-custom entries (`:octocat:`, `:shipit:`, etc.) that have no unicode
codepoint.

### Regenerating shortcode data

```bash
curl -s https://api.github.com/emojis > emojis.json
go generate ./...
```

## Related

- [github.com/client9/goldmark-demoji](https://github.com/client9/goldmark-demoji) —
  goldmark extension wrapper with functional options

## License

[MIT](/LICENSE)

