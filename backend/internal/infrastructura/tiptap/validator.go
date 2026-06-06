package tiptap

import (
	"encoding/json"
	"fmt"
)

var allowedNodeTypes = map[string]bool{
	"doc": true, "paragraph": true, "text": true,
	"heading": true, "bulletList": true, "orderedList": true,
	"listItem": true, "blockquote": true, "codeBlock": true,
	"hardBreak": true, "horizontalRule": true, "image": true,
	"table": true, "tableRow": true, "tableCell": true, "tableHeader": true,
	"mention": true, "taskList": true, "taskItem": true,
}

var allowedMarkTypes = map[string]bool{
	"bold": true, "italic": true, "underline": true,
	"strike": true, "code": true, "link": true,
	"highlight": true, "textStyle": true,
}

// Validate checks that raw JSON is a valid TipTap document with known node/mark types.
func Validate(raw json.RawMessage) error {
	if len(raw) == 0 {
		return nil
	}
	var doc Document
	if err := json.Unmarshal(raw, &doc); err != nil {
		return fmt.Errorf("tiptap: invalid JSON: %w", err)
	}
	if doc.Type != "doc" {
		return fmt.Errorf("tiptap: root type must be \"doc\", got %q", doc.Type)
	}
	return validateNodes(doc.Content)
}

func validateNodes(nodes []Node) error {
	for _, n := range nodes {
		if !allowedNodeTypes[n.Type] {
			return fmt.Errorf("tiptap: unknown node type %q", n.Type)
		}
		for _, m := range n.Marks {
			if !allowedMarkTypes[m.Type] {
				return fmt.Errorf("tiptap: unknown mark type %q", m.Type)
			}
		}
		if err := validateNodes(n.Content); err != nil {
			return err
		}
	}
	return nil
}
