// Converts GitHub-Flavored Markdown (+ Obsidian callouts) to TipTap JSON.
// Handles: headings, bold/italic/code/strike, links, images, bullet/ordered/task
// lists, code blocks, blockquotes, horizontal rules, tables.

export function markdownToTiptap(md) {
  const lines = md.replace(/\r\n/g, '\n').split('\n');
  const nodes = [];
  let i = 0;

  while (i < lines.length) {
    const line = lines[i];

    // ── Fenced code block ─────────────────────────────────────────
    if (/^```/.test(line)) {
      const lang = line.slice(3).trim() || null;
      const codeLines = [];
      i++;
      while (i < lines.length && !/^```/.test(lines[i])) {
        codeLines.push(lines[i]);
        i++;
      }
      nodes.push({
        type: 'codeBlock',
        attrs: { language: lang },
        content: [{ type: 'text', text: codeLines.join('\n') }],
      });
      i++; // skip closing ```
      continue;
    }

    // ── Heading ────────────────────────────────────────────────────
    const hm = line.match(/^(#{1,6})\s+(.*)/);
    if (hm) {
      const level = Math.min(hm[1].length, 3);
      nodes.push({ type: 'heading', attrs: { level }, content: parseInline(hm[2]) });
      i++;
      continue;
    }

    // ── Horizontal rule ────────────────────────────────────────────
    if (/^(-{3,}|\*{3,}|_{3,})$/.test(line.trim())) {
      nodes.push({ type: 'horizontalRule' });
      i++;
      continue;
    }

    // ── Blockquote (including Obsidian callouts) ────────────────────
    if (/^>\s?/.test(line)) {
      const bqLines = [];
      while (i < lines.length && /^>\s?/.test(lines[i])) {
        bqLines.push(lines[i].replace(/^>\s?/, ''));
        i++;
      }
      // Strip Obsidian callout marker e.g. [!NOTE] on first line
      const filtered = bqLines.filter((l, idx) => !(idx === 0 && /^\[![A-Z]+\]/.test(l)));
      const inner = filtered.join('\n').trim();
      if (inner) {
        nodes.push({
          type: 'blockquote',
          content: [{ type: 'paragraph', content: parseInline(inner) }],
        });
      }
      continue;
    }

    // ── Task list ──────────────────────────────────────────────────
    if (/^[-*+]\s+\[[ xX]\]/.test(line)) {
      const taskItems = [];
      while (i < lines.length && /^[-*+]\s+\[[ xX]\]/.test(lines[i])) {
        const checked = /^[-*+]\s+\[[xX]\]/.test(lines[i]);
        const text = lines[i].replace(/^[-*+]\s+\[[ xX]\]\s*/, '');
        taskItems.push({
          type: 'taskItem',
          attrs: { checked },
          content: [{ type: 'paragraph', content: parseInline(text) }],
        });
        i++;
      }
      nodes.push({ type: 'taskList', content: taskItems });
      continue;
    }

    // ── Bullet list ────────────────────────────────────────────────
    if (/^[-*+]\s/.test(line)) {
      const items = [];
      while (i < lines.length && /^[-*+]\s/.test(lines[i])) {
        const text = lines[i].replace(/^[-*+]\s+/, '');
        items.push({
          type: 'listItem',
          content: [{ type: 'paragraph', content: parseInline(text) }],
        });
        i++;
      }
      nodes.push({ type: 'bulletList', content: items });
      continue;
    }

    // ── Ordered list ───────────────────────────────────────────────
    if (/^\d+[.)]\s/.test(line)) {
      const items = [];
      while (i < lines.length && /^\d+[.)]\s/.test(lines[i])) {
        const text = lines[i].replace(/^\d+[.)]\s+/, '');
        items.push({
          type: 'listItem',
          content: [{ type: 'paragraph', content: parseInline(text) }],
        });
        i++;
      }
      nodes.push({ type: 'orderedList', attrs: { start: 1 }, content: items });
      continue;
    }

    // ── Table ──────────────────────────────────────────────────────
    if (/^\|/.test(line)) {
      const tableLines = [];
      while (i < lines.length && /^\|/.test(lines[i])) {
        tableLines.push(lines[i]);
        i++;
      }
      const table = parseTable(tableLines);
      if (table) nodes.push(table);
      continue;
    }

    // ── Empty line ─────────────────────────────────────────────────
    if (line.trim() === '') {
      i++;
      continue;
    }

    // ── Paragraph (collects lines until a block boundary) ──────────
    const paraLines = [];
    while (
      i < lines.length &&
      lines[i].trim() !== '' &&
      !/^#{1,6}\s/.test(lines[i]) &&
      !/^[-*+]\s/.test(lines[i]) &&
      !/^\d+[.)]\s/.test(lines[i]) &&
      !/^```/.test(lines[i]) &&
      !/^>\s?/.test(lines[i]) &&
      !/^\|/.test(lines[i]) &&
      !/^(-{3,}|\*{3,}|_{3,})$/.test(lines[i].trim())
    ) {
      paraLines.push(lines[i]);
      i++;
    }
    if (paraLines.length > 0) {
      nodes.push({ type: 'paragraph', content: parseInline(paraLines.join('\n')) });
    }
  }

  return {
    type: 'doc',
    content: nodes.length > 0 ? nodes : [{ type: 'paragraph' }],
  };
}

// ── Table parser ──────────────────────────────────────────────────────────────
function parseTable(lines) {
  if (lines.length < 2) return null;
  const splitRow = (l) =>
    l.split('|')
      .slice(1, -1)
      .map((c) => c.trim());

  const headers = splitRow(lines[0]);
  const bodyLines = lines.slice(2); // skip separator row

  const headerRow = {
    type: 'tableRow',
    content: headers.map((cell) => ({
      type: 'tableHeader',
      attrs: { colspan: 1, rowspan: 1, colwidth: null },
      content: [{ type: 'paragraph', content: parseInline(cell) }],
    })),
  };

  const bodyRows = bodyLines.map((line) => ({
    type: 'tableRow',
    content: splitRow(line).map((cell) => ({
      type: 'tableCell',
      attrs: { colspan: 1, rowspan: 1, colwidth: null },
      content: [{ type: 'paragraph', content: parseInline(cell) }],
    })),
  }));

  return { type: 'table', content: [headerRow, ...bodyRows] };
}

// ── Inline parser ─────────────────────────────────────────────────────────────
// Returns an array of TipTap inline nodes (text with optional marks, or image).
function parseInline(text) {
  if (!text) return [{ type: 'text', text: '' }];

  const patterns = [
    // Image — must come before link
    {
      re: /!\[([^\]]*)\]\(([^)]+)\)/,
      node: (m) => ({ type: 'image', attrs: { src: m[2], alt: m[1], title: null } }),
    },
    // Link
    {
      re: /\[([^\]]+)\]\(([^)]+)\)/,
      node: (m) => ({
        type: 'text',
        text: m[1],
        marks: [{ type: 'link', attrs: { href: m[2], target: '_blank', rel: 'noopener noreferrer' } }],
      }),
    },
    // Bold + italic
    { re: /\*\*\*(.+?)\*\*\*/, node: (m) => ({ type: 'text', text: m[1], marks: [{ type: 'bold' }, { type: 'italic' }] }) },
    // Bold (**text** or __text__)
    { re: /\*\*(.+?)\*\*|__(.+?)__/, node: (m) => ({ type: 'text', text: m[1] ?? m[2], marks: [{ type: 'bold' }] }) },
    // Italic (*text* or _text_) — single chars, not inside words for _
    { re: /\*([^*]+?)\*/, node: (m) => ({ type: 'text', text: m[1], marks: [{ type: 'italic' }] }) },
    // Strikethrough
    { re: /~~(.+?)~~/, node: (m) => ({ type: 'text', text: m[1], marks: [{ type: 'strike' }] }) },
    // Inline code
    { re: /`([^`]+?)`/, node: (m) => ({ type: 'text', text: m[1], marks: [{ type: 'code' }] }) },
  ];

  const result = [];
  let remaining = text;

  while (remaining.length > 0) {
    let earliest = null;
    let earliestIdx = Infinity;

    for (const p of patterns) {
      const m = p.re.exec(remaining);
      if (m && m.index < earliestIdx) {
        earliest = { match: m, node: p.node };
        earliestIdx = m.index;
      }
    }

    if (!earliest) {
      result.push({ type: 'text', text: remaining });
      break;
    }

    if (earliestIdx > 0) {
      result.push({ type: 'text', text: remaining.slice(0, earliestIdx) });
    }

    result.push(earliest.node(earliest.match));
    remaining = remaining.slice(earliestIdx + earliest.match[0].length);
  }

  return result.length > 0 ? result : [{ type: 'text', text: '' }];
}
