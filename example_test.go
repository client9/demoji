package demoji_test

import (
	"fmt"
	"maps"

	"github.com/client9/demoji"
)

func ExampleReplace() {
	fmt.Println(demoji.Replace("Hello 🙂 world"))
	// Output: Hello :-) world
}

func ExampleReplaceBytes() {
	fmt.Println(string(demoji.ReplaceBytes([]byte("I ❤️ Go"))))
	// Output: I <3 Go
}

func ExampleNew_unicodeToShortcode() {
	r := demoji.New(demoji.Config{
		From: demoji.FormatUnicode,
		To:   demoji.FormatShortcode,
	})
	fmt.Println(r.Replace("Hello 🙂 world"))
	// Output: Hello :slightly_smiling_face: world
}

func ExampleNew_emoticonToUnicode() {
	r := demoji.New(demoji.Config{From: demoji.FormatEmoticon})
	fmt.Println(r.Replace("Hello :-) world"))
	// Output: Hello 🙂 world
}

func ExampleNew_shortcodeToUnicode() {
	r := demoji.New(demoji.Config{From: demoji.FormatShortcode})
	fmt.Println(r.Replace("I :heart: Go"))
	// Output: I ❤ Go
}

func ExampleNew_combinedToUnicode() {
	r := demoji.New(demoji.Config{From: demoji.FormatEmoticon | demoji.FormatShortcode})
	fmt.Println(r.Replace(":-) and :grinning: are both happy"))
	// Output: 🙂 and 😀 are both happy
}

func ExampleNew_customEmoticon() {
	emoticons := maps.Clone(demoji.DefaultEmoticons())
	emoticons["^^"] = "😊"
	r := demoji.New(demoji.Config{
		From:      demoji.FormatEmoticon,
		Emoticons: emoticons,
	})
	fmt.Println(r.Replace("Hello ^^ world"))
	// Output: Hello 😊 world
}

func ExampleNew_excludeEmoticon() {
	// Remove ":-)" to shift the canonical to ":)".
	emoticons := maps.Clone(demoji.DefaultEmoticons())
	delete(emoticons, ":-)")
	r := demoji.New(demoji.Config{
		From:      demoji.FormatUnicode,
		To:        demoji.FormatEmoticon,
		Emoticons: emoticons,
	})
	fmt.Println(r.Replace("Hello 🙂 world"))
	// Output: Hello :) world
}
