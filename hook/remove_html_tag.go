package hook

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

// https://core.telegram.org/bots/api#sendmessage
func RemoveHTMLTag(content string) string {
	p := bluemonday.NewPolicy()

	p.AllowStandardURLs()
	// <b>bold</b>
	p.AllowElements("b")
	// <strong>bold</strong>
	p.AllowElements("strong")
	// <i>italic</i>
	p.AllowElements("i")
	// <i>italic</i>
	p.AllowElements("em")
	// <u>underline</u>
	p.AllowElements("u")
	// <ins>underline</ins>
	p.AllowElements("ins")
	// <s>strikethrough</s>
	p.AllowElements("s")
	// <strike>strikethrough</strike>
	p.AllowElements("strike")
	// <del>strikethrough</del>
	p.AllowElements("del")
	// <span class="tg-spoiler">spoiler</span>
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^tg-spoiler$`)).OnElements("span")
	// <a href="http://www.example.com/">inline URL</a>
	p.AllowAttrs("href").OnElements("a")
	p.RequireNoFollowOnLinks(false)
	p.AllowURLSchemes("http", "https", "tg")
	// <tg-emoji emoji-id="5368324170671202286">👍</tg-emoji>
	p.AllowAttrs("emoji-id").Matching(regexp.MustCompile(`^\d+$`)).OnElements("tg-emoji")
	// <code>inline fixed-width code</code>
	p.AllowElements("code")
	// <pre>pre-formatted fixed-width code block</pre>
	p.AllowElements("pre")
	// <pre><code class="language-python">pre-formatted fixed-width code block written in the Python programming language</code></pre>
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^language-[\w-]+$`)).OnElements("code")
	// <blockquote>Block quotation started\nBlock quotation continued\nThe last line of the block quotation</blockquote>
	p.AllowElements("blockquote")

	// <tg-spoiler>spoiler</tg-spoiler>
	// 这条测不过，因为 <tg-spoiler>spoiler</tg-spoiler> 是一个自定义的标签，不是标准的 HTML 标签
	// p.addDefaultElementsWithoutAttrs 中没有定义
	// p.AllowElements("tg-spoiler")
	// <blockquote expandable>Expandable block quotation started\nExpandable block quotation continued\nExpandable block quotation continued\nHidden by default part of the block quotation started\nExpandable block quotation continued\nThe last line of the block quotation</blockquote>
	// 这条测不过，因为 <blockquote expandable> 是一个自定义的属性，不是标准的 HTML 属性
	// p.AllowAttrs("expandable").OnElements("blockquote")
	// 不过也没有影响，正常邮件中不会有上述标签、属性

	return p.Sanitize(content)
}
