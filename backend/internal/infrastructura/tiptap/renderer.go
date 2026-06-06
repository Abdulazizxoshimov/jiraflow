package tiptap

import (
	"fmt"
	"html"
	"strings"
)

// RenderHTML converts a TipTap Document JSON to safe HTML.
// Frontend ushbu HTMLni to'g'ridan-to'g'ri ko'rsatishi mumkin.
func (d *Document) RenderHTML() string {
	var sb strings.Builder
	renderNodes(&sb, d.Content)
	return sb.String()
}

// RenderMarkdown converts a TipTap Document to GitHub-flavored Markdown.
func (d *Document) RenderMarkdown() string {
	var sb strings.Builder
	renderMarkdownNodes(&sb, d.Content, 0)
	return strings.TrimSpace(sb.String())
}

func renderMarkdownNodes(sb *strings.Builder, nodes []Node, listDepth int) {
	for _, n := range nodes {
		renderMarkdownNode(sb, n, listDepth)
	}
}

func renderMarkdownNode(sb *strings.Builder, n Node, listDepth int) {
	switch n.Type {
	case "text":
		renderMarkdownText(sb, n)
	case "paragraph":
		renderMarkdownNodes(sb, n.Content, listDepth)
		sb.WriteString("\n\n")
	case "heading":
		level := intAttr(n, "level", 1)
		sb.WriteString(strings.Repeat("#", level) + " ")
		renderMarkdownNodes(sb, n.Content, listDepth)
		sb.WriteString("\n\n")
	case "bulletList":
		renderMarkdownNodes(sb, n.Content, listDepth+1)
		sb.WriteString("\n")
	case "orderedList":
		renderMarkdownNodes(sb, n.Content, listDepth+1)
		sb.WriteString("\n")
	case "listItem":
		sb.WriteString(strings.Repeat("  ", listDepth-1) + "- ")
		renderMarkdownNodes(sb, n.Content, listDepth)
	case "taskList":
		renderMarkdownNodes(sb, n.Content, listDepth+1)
		sb.WriteString("\n")
	case "taskItem":
		checked := boolAttr(n, "checked")
		if checked {
			sb.WriteString(strings.Repeat("  ", listDepth-1) + "- [x] ")
		} else {
			sb.WriteString(strings.Repeat("  ", listDepth-1) + "- [ ] ")
		}
		renderMarkdownNodes(sb, n.Content, listDepth)
	case "blockquote":
		var inner strings.Builder
		renderMarkdownNodes(&inner, n.Content, listDepth)
		for _, line := range strings.Split(strings.TrimRight(inner.String(), "\n"), "\n") {
			sb.WriteString("> " + line + "\n")
		}
		sb.WriteString("\n")
	case "codeBlock":
		lang := strAttr(n, "language")
		sb.WriteString("```" + lang + "\n")
		renderMarkdownNodes(sb, n.Content, listDepth)
		sb.WriteString("```\n\n")
	case "hardBreak":
		sb.WriteString("  \n")
	case "horizontalRule":
		sb.WriteString("---\n\n")
	case "image":
		src := strAttr(n, "src")
		alt := strAttr(n, "alt")
		sb.WriteString(fmt.Sprintf("![%s](%s)\n\n", alt, src))
	case "mention":
		sb.WriteString("@" + strAttr(n, "label"))
	case "table":
		renderMarkdownTable(sb, n)
	}
}

func renderMarkdownText(sb *strings.Builder, n Node) {
	text := n.Text
	bold, italic, code, strike := false, false, false, false
	var link string
	for _, m := range n.Marks {
		switch m.Type {
		case "bold":
			bold = true
		case "italic":
			italic = true
		case "code":
			code = true
		case "strike":
			strike = true
		case "link":
			link = strAttrFromMap(m.Attrs, "href")
		}
	}
	if code {
		sb.WriteString("`" + text + "`")
		return
	}
	if bold {
		text = "**" + text + "**"
	}
	if italic {
		text = "_" + text + "_"
	}
	if strike {
		text = "~~" + text + "~~"
	}
	if link != "" {
		text = "[" + text + "](" + link + ")"
	}
	sb.WriteString(text)
}

func renderMarkdownTable(sb *strings.Builder, n Node) {
	if len(n.Content) == 0 {
		return
	}
	for i, row := range n.Content {
		var cells []string
		for _, cell := range row.Content {
			var cellSB strings.Builder
			renderMarkdownNodes(&cellSB, cell.Content, 0)
			cells = append(cells, strings.TrimSpace(cellSB.String()))
		}
		sb.WriteString("| " + strings.Join(cells, " | ") + " |\n")
		if i == 0 {
			sep := make([]string, len(cells))
			for j := range sep {
				sep[j] = "---"
			}
			sb.WriteString("| " + strings.Join(sep, " | ") + " |\n")
		}
	}
	sb.WriteString("\n")
}

// PlainText extracts all text from the document — qidiruv indekslash va
// email preview uchun ishlatiladi.
func (d *Document) PlainText() string {
	var sb strings.Builder
	extractText(&sb, d.Content)
	return strings.TrimSpace(sb.String())
}

// Excerpt returns up to maxLen runes of plain text.
func (d *Document) Excerpt(maxLen int) string {
	plain := d.PlainText()
	runes := []rune(plain)
	if len(runes) <= maxLen {
		return plain
	}
	return string(runes[:maxLen]) + "…"
}

// ─── Node renderer ────────────────────────────────────────────────────────────

func renderNodes(sb *strings.Builder, nodes []Node) {
	for _, n := range nodes {
		renderNode(sb, n)
	}
}

func renderNode(sb *strings.Builder, n Node) {
	switch n.Type {
	case "text":
		renderTextNode(sb, n)

	case "paragraph":
		sb.WriteString("<p>")
		renderNodes(sb, n.Content)
		sb.WriteString("</p>\n")

	case "heading":
		level := intAttr(n, "level", 1)
		if level < 1 || level > 6 {
			level = 1
		}
		tag := fmt.Sprintf("h%d", level)
		sb.WriteString("<" + tag + ">")
		renderNodes(sb, n.Content)
		sb.WriteString("</" + tag + ">\n")

	case "bulletList":
		sb.WriteString("<ul>\n")
		renderNodes(sb, n.Content)
		sb.WriteString("</ul>\n")

	case "orderedList":
		start := intAttr(n, "start", 1)
		if start != 1 {
			sb.WriteString(fmt.Sprintf(`<ol start="%d">`+"\n", start))
		} else {
			sb.WriteString("<ol>\n")
		}
		renderNodes(sb, n.Content)
		sb.WriteString("</ol>\n")

	case "listItem":
		sb.WriteString("<li>")
		renderNodes(sb, n.Content)
		sb.WriteString("</li>\n")

	case "taskList":
		sb.WriteString(`<ul class="task-list">` + "\n")
		renderNodes(sb, n.Content)
		sb.WriteString("</ul>\n")

	case "taskItem":
		checked := boolAttr(n, "checked")
		if checked {
			sb.WriteString(`<li class="task-item"><input type="checkbox" checked disabled> `)
		} else {
			sb.WriteString(`<li class="task-item"><input type="checkbox" disabled> `)
		}
		renderNodes(sb, n.Content)
		sb.WriteString("</li>\n")

	case "blockquote":
		sb.WriteString("<blockquote>\n")
		renderNodes(sb, n.Content)
		sb.WriteString("</blockquote>\n")

	case "codeBlock":
		lang := strAttr(n, "language")
		if lang != "" {
			sb.WriteString(fmt.Sprintf("<pre><code class=\"language-%s\">", html.EscapeString(lang)))
		} else {
			sb.WriteString("<pre><code>")
		}
		renderNodes(sb, n.Content)
		sb.WriteString("</code></pre>\n")

	case "hardBreak":
		sb.WriteString("<br>")

	case "horizontalRule":
		sb.WriteString("<hr>\n")

	case "image":
		src := html.EscapeString(strAttr(n, "src"))
		alt := html.EscapeString(strAttr(n, "alt"))
		title := strAttr(n, "title")
		if title != "" {
			sb.WriteString(fmt.Sprintf(`<img src="%s" alt="%s" title="%s">`, src, alt, html.EscapeString(title)))
		} else {
			sb.WriteString(fmt.Sprintf(`<img src="%s" alt="%s">`, src, alt))
		}

	case "table":
		sb.WriteString(`<table class="tiptap-table">` + "\n")
		renderNodes(sb, n.Content)
		sb.WriteString("</table>\n")

	case "tableRow":
		sb.WriteString("<tr>\n")
		renderNodes(sb, n.Content)
		sb.WriteString("</tr>\n")

	case "tableHeader":
		colspan := intAttr(n, "colspan", 1)
		rowspan := intAttr(n, "rowspan", 1)
		attrs := cellAttrs(colspan, rowspan)
		sb.WriteString("<th" + attrs + ">")
		renderNodes(sb, n.Content)
		sb.WriteString("</th>\n")

	case "tableCell":
		colspan := intAttr(n, "colspan", 1)
		rowspan := intAttr(n, "rowspan", 1)
		attrs := cellAttrs(colspan, rowspan)
		sb.WriteString("<td" + attrs + ">")
		renderNodes(sb, n.Content)
		sb.WriteString("</td>\n")

	case "mention":
		id := html.EscapeString(strAttr(n, "id"))
		label := html.EscapeString(strAttr(n, "label"))
		sb.WriteString(fmt.Sprintf(`<span class="mention" data-id="%s">@%s</span>`, id, label))
	}
}

// ─── Text + marks ─────────────────────────────────────────────────────────────

func renderTextNode(sb *strings.Builder, n Node) {
	text := html.EscapeString(n.Text)
	if len(n.Marks) == 0 {
		sb.WriteString(text)
		return
	}

	// Marklarni ochish (original tartibda)
	for _, m := range n.Marks {
		sb.WriteString(openMark(m))
	}
	sb.WriteString(text)
	// Marklarni yopish (teskari tartibda)
	for i := len(n.Marks) - 1; i >= 0; i-- {
		sb.WriteString(closeMark(n.Marks[i]))
	}
}

func openMark(m Mark) string {
	switch m.Type {
	case "bold":
		return "<strong>"
	case "italic":
		return "<em>"
	case "underline":
		return "<u>"
	case "strike":
		return "<s>"
	case "code":
		return "<code>"
	case "highlight":
		if color := strAttrFromMap(m.Attrs, "color"); color != "" {
			return fmt.Sprintf(`<mark style="background-color:%s">`, html.EscapeString(color))
		}
		return "<mark>"
	case "link":
		href := html.EscapeString(strAttrFromMap(m.Attrs, "href"))
		target := strAttrFromMap(m.Attrs, "target")
		if target == "" {
			target = "_blank"
		}
		return fmt.Sprintf(`<a href="%s" target="%s" rel="noopener noreferrer">`,
			href, html.EscapeString(target))
	case "textStyle":
		var style strings.Builder
		if color := strAttrFromMap(m.Attrs, "color"); color != "" {
			style.WriteString(fmt.Sprintf("color:%s;", html.EscapeString(color)))
		}
		if fontSize := strAttrFromMap(m.Attrs, "fontSize"); fontSize != "" {
			style.WriteString(fmt.Sprintf("font-size:%s;", html.EscapeString(fontSize)))
		}
		if style.Len() > 0 {
			return fmt.Sprintf(`<span style="%s">`, style.String())
		}
		return ""
	}
	return ""
}

func closeMark(m Mark) string {
	switch m.Type {
	case "bold":
		return "</strong>"
	case "italic":
		return "</em>"
	case "underline":
		return "</u>"
	case "strike":
		return "</s>"
	case "code":
		return "</code>"
	case "highlight":
		return "</mark>"
	case "link":
		return "</a>"
	case "textStyle":
		if strAttrFromMap(m.Attrs, "color") != "" || strAttrFromMap(m.Attrs, "fontSize") != "" {
			return "</span>"
		}
		return ""
	}
	return ""
}

// ─── PlainText extractor ──────────────────────────────────────────────────────

func extractText(sb *strings.Builder, nodes []Node) {
	for _, n := range nodes {
		switch n.Type {
		case "text":
			sb.WriteString(n.Text)
		case "hardBreak":
			sb.WriteString("\n")
		case "paragraph", "heading", "blockquote", "listItem", "taskItem":
			extractText(sb, n.Content)
			sb.WriteString("\n")
		case "codeBlock":
			extractText(sb, n.Content)
			sb.WriteString("\n")
		case "horizontalRule":
			sb.WriteString("\n---\n")
		case "image":
			if alt := strAttr(n, "alt"); alt != "" {
				sb.WriteString("[" + alt + "]")
			}
		case "mention":
			sb.WriteString("@" + strAttr(n, "label"))
		default:
			extractText(sb, n.Content)
		}
	}
}

// ─── Attribute helpers ────────────────────────────────────────────────────────

func strAttr(n Node, key string) string {
	return strAttrFromMap(n.Attrs, key)
}

func strAttrFromMap(attrs map[string]any, key string) string {
	if attrs == nil {
		return ""
	}
	v, _ := attrs[key].(string)
	return v
}

func intAttr(n Node, key string, def int) int {
	if n.Attrs == nil {
		return def
	}
	switch v := n.Attrs[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	}
	return def
}

func boolAttr(n Node, key string) bool {
	if n.Attrs == nil {
		return false
	}
	v, _ := n.Attrs[key].(bool)
	return v
}

func cellAttrs(colspan, rowspan int) string {
	var sb strings.Builder
	if colspan > 1 {
		sb.WriteString(fmt.Sprintf(` colspan="%d"`, colspan))
	}
	if rowspan > 1 {
		sb.WriteString(fmt.Sprintf(` rowspan="%d"`, rowspan))
	}
	return sb.String()
}
