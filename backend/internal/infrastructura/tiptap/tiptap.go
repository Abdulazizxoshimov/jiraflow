package tiptap

import "encoding/json"

// Document represents a TipTap/ProseMirror JSON document.
type Document struct {
	Type    string `json:"type"`
	Content []Node `json:"content,omitempty"`
}

// Node is a single ProseMirror node.
type Node struct {
	Type    string         `json:"type"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Content []Node         `json:"content,omitempty"`
	Marks   []Mark         `json:"marks,omitempty"`
	Text    string         `json:"text,omitempty"`
}

// Mark is inline formatting applied to a text node.
type Mark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

// Parse deserializes raw JSON into a Document.
func Parse(raw []byte) (*Document, error) {
	var doc Document
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// Bytes serializes the document back to JSON.
func (d *Document) Bytes() ([]byte, error) {
	return json.Marshal(d)
}

// Empty returns true when the document has no content nodes.
func (d *Document) Empty() bool {
	return len(d.Content) == 0
}

// ExtractMentionIDs walks a TipTap document stored as map[string]any
// and returns all unique user IDs from mention nodes.
func ExtractMentionIDs(raw map[string]any) []string {
	seen := make(map[string]struct{})
	var ids []string
	walkMap(raw, seen, &ids)
	return ids
}

func walkMap(node map[string]any, seen map[string]struct{}, ids *[]string) {
	if t, _ := node["type"].(string); t == "mention" {
		if attrs, ok := node["attrs"].(map[string]any); ok {
			if id, _ := attrs["id"].(string); id != "" {
				if _, dup := seen[id]; !dup {
					seen[id] = struct{}{}
					*ids = append(*ids, id)
				}
			}
		}
	}
	if content, ok := node["content"].([]any); ok {
		for _, item := range content {
			if child, ok := item.(map[string]any); ok {
				walkMap(child, seen, ids)
			}
		}
	}
}
