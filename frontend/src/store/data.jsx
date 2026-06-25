// data.jsx — Static UI meta (colors, type icons, priority icons).

export const COLORS = ["#6366F1","#06B6D4","#10B981","#F59E0B","#EF4444","#8B5CF6","#EC4899","#14B8A6","#F97316","#3B82F6"];

export const TYPE_META = {
  Bug:     { tone: "danger",  icon: "bug",   color: "#EF4444" },
  Task:    { tone: "info",    icon: "task",  color: "#3B82F6" },
  Story:   { tone: "purple",  icon: "story", color: "#8B5CF6" },
  Epic:    { tone: "orange",  icon: "epic",  color: "#F97316" },
  Subtask: { tone: "muted",   icon: "task",  color: "#64748B" },
};

export const PRIORITY_META = {
  Critical: { tone: "danger",  icon: "prHigh", color: "#DC2626" },
  High:     { tone: "warning", icon: "prHigh", color: "#EA580C" },
  Medium:   { tone: "info",    icon: "prMed",  color: "#3B82F6" },
  Low:      { tone: "muted",   icon: "prLow",  color: "#64748B" },
};

export const STATUS_META = {
  Backlog:       { tone: "muted"   },
  "To Do":       { tone: "muted"   },
  Todo:          { tone: "muted"   },
  "In Progress": { tone: "info"    },
  "In Review":   { tone: "purple"  },
  Done:          { tone: "success" },
  Open:          { tone: "danger"  },
  Triaged:       { tone: "warning" },
  Fixing:        { tone: "info"    },
  Verifying:     { tone: "purple"  },
  Closed:        { tone: "success" },
};
