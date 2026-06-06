package page

import "github.com/jira-backend/jiraflow-backend/internal/entity"

// buildTree converts a flat slice into a nested tree.
// Root nodes have nil ParentID; children are attached under their parent.
func buildTree(nodes []*entity.PageTree) []*entity.PageTree {
	index := make(map[string]*entity.PageTree, len(nodes))
	for _, n := range nodes {
		index[n.ID] = n
	}

	var roots []*entity.PageTree
	for _, n := range nodes {
		if n.ParentID == nil {
			roots = append(roots, n)
		} else {
			if parent, ok := index[*n.ParentID]; ok {
				parent.Children = append(parent.Children, *n)
			} else {
				roots = append(roots, n)
			}
		}
	}
	return roots
}
