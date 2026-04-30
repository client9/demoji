# demoji

A [goldmark](https://github.com/yuin/goldmark) extension that converts between unicode emoji, classic ASCII emoticons, and GitHub-style shortcodes — in either direction.

```
😀  ↔  :-D          (unicode ↔ emoticon)
🙂  ↔  :slightly_smiling_face:   (unicode ↔ shortcode)
:-) ↔  :slightly_smiling_face:   (emoticon ↔ shortcode, via two passes)
```

## Install

```bash
go get github.com/client9/demoji
```

Requires Go 1.22+ and goldmark v1.7+.

## Quick start

```go
import (
    "github.com/client9/demoji"
    "github.com/yuin/goldmark"
)

// unicode emoji → ASCII emoticons (default)
md := goldmark.New(goldmark.WithExtensions(demoji.New()))

// ASCII emoticons → unicode emoji
md := goldmark.New(goldmark.WithExtensions(
    demoji.New(demoji.WithFrom(demoji.FormatEmoticon)),
))

// GitHub shortcodes → unicode emoji
md := goldmark.New(goldmark.WithExtensions(
    demoji.New(demoji.WithFrom(demoji.FormatShortcode)),
))

// unicode emoji → GitHub shortcodes
md := goldmark.New(goldmark.WithExtensions(
    demoji.New(
        demoji.WithFrom(demoji.FormatUnicode),
        demoji.WithTo(demoji.FormatShortcode),
    ),
))
```

## Formats

| Constant | Example | Description |
|----------|---------|-------------|
| `FormatUnicode` | `😀` | Actual unicode emoji codepoints |
| `FormatEmoticon` | `:-D` | Classic ASCII emoticons |
| `FormatShortcode` | `:grinning:` | GitHub / Gemoji shortcode names |

`WithFrom` accepts a bitfield — combine formats to convert multiple sources to unicode in one pass:

```go
// emoticons AND shortcodes → unicode in a single AST walk
md := goldmark.New(goldmark.WithExtensions(
    demoji.New(
        demoji.WithFrom(demoji.FormatEmoticon | demoji.FormatShortcode),
    ),
))
```

Input `":-) and :grinning: both work"` → `"🙂 and 😀 both work"`.

## Built-in mappings

**Emoticons** (~40 entries): the classic set — `:-)` `:)` `:-D` `:-(`  `;-)` `:-P` `:-O` `B-)` `O:-)` `>:(` `:'(` `<3` `</3` `XD` `-_-` and their variants.

**Shortcodes** (1913 entries): the full [GitHub emoji](https://api.github.com/emojis) set, minus 23 GitHub-custom images (`:octocat:`, `:shipit:`, etc.) that have no unicode codepoint.

Content inside code spans and fenced code blocks is never modified.

## Configuration

### Add mappings

```go
demoji.New(
    demoji.WithFrom(demoji.FormatEmoticon),
    demoji.WithAdditional(demoji.FormatEmoticon, map[string]string{
        "^^":  "😊",
        "o/":  "👋",
    }),
)
```

### Remove mappings

```go
// stop converting XD and xD
demoji.New(
    demoji.WithExclude(demoji.FormatEmoticon, "XD", "xD"),
)

// stop converting a shortcode
demoji.New(
    demoji.WithFrom(demoji.FormatShortcode),
    demoji.WithExclude(demoji.FormatShortcode, ":cry:"),
)
```

### Shift the canonical emoticon

In unicode→emoticon mode each emoji has one canonical text form (the first active entry in the builtin list). Excluding the current canonical promotes the next entry:

```go
// default: 🙂 → :-)
// after exclude: 🙂 → :)
demoji.New(
    demoji.WithExclude(demoji.FormatEmoticon, ":-)"),
)
```

## Regenerating shortcode data

The shortcode table is generated from a local copy of the GitHub emoji API response:

```bash
# refresh the source data
curl -s https://api.github.com/emojis > emojis.json

# regenerate shortcodes_gen.go
go generate ./...
```

## How it works

Conversions that search for text patterns (emoticon→unicode, shortcode→unicode) use goldmark's inline parser and AST transformer APIs so that code spans and fenced blocks are automatically excluded.

Shortcode→unicode uses a dedicated **inline parser** triggered on `:` rather than an AST transformer. This is necessary because goldmark's emphasis tokenizer splits text at `_` boundaries before AST transformers run, which would break shortcode names like `:slightly_smiling_face:`.

When `WithFrom` specifies multiple formats, both mechanisms run cooperatively in a single parse — the inline parser handles shortcodes during tokenization, the AST transformer handles emoticons afterward.

## License

MIT
