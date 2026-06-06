CREATE TABLE blog_posts (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id     UUID         NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    title        VARCHAR(500) NOT NULL,
    body         TEXT         NOT NULL DEFAULT '',
    author_id    UUID         NOT NULL REFERENCES users(id),
    is_published BOOLEAN      NOT NULL DEFAULT FALSE,
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_blog_posts_space  ON blog_posts (space_id, published_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_blog_posts_author ON blog_posts (author_id)                  WHERE deleted_at IS NULL;

CREATE TRIGGER trg_blog_posts_updated_at
    BEFORE UPDATE ON blog_posts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
