-- Issue ↔ Page ko'prik jadvali
CREATE TABLE IF NOT EXISTS issue_page_links (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    page_id    UUID NOT NULL REFERENCES pages(id)  ON DELETE CASCADE,
    linked_by  UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (issue_id, page_id)
);

CREATE INDEX idx_issue_page_links_issue ON issue_page_links(issue_id);
CREATE INDEX idx_issue_page_links_page  ON issue_page_links(page_id);

-- Unified activity feed jadvali
CREATE TABLE IF NOT EXISTS activity_feed (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id     UUID NOT NULL REFERENCES users(id),
    action       TEXT NOT NULL,    -- created | updated | deleted | commented | transitioned | linked | mentioned
    entity_type  TEXT NOT NULL,    -- issue | page | comment | sprint | space | project
    entity_id    UUID NOT NULL,
    entity_title TEXT NOT NULL,
    project_id   UUID REFERENCES projects(id) ON DELETE SET NULL,
    space_id     UUID REFERENCES spaces(id)   ON DELETE SET NULL,
    meta         JSONB,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_feed_actor      ON activity_feed(actor_id);
CREATE INDEX idx_activity_feed_project    ON activity_feed(project_id);
CREATE INDEX idx_activity_feed_space      ON activity_feed(space_id);
CREATE INDEX idx_activity_feed_entity     ON activity_feed(entity_type, entity_id);
CREATE INDEX idx_activity_feed_created_at ON activity_feed(created_at DESC);

-- Project-Space avto-bog'lash uchun spaces jadvaliga project_id indeks (agar yo'q bo'lsa)
CREATE INDEX IF NOT EXISTS idx_spaces_project_id ON spaces(project_id);
