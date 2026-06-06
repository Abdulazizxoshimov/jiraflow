INSERT INTO users (
    id,
    email,
    password_hash,
    full_name,
    role,
    is_active
) VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'admin@jiraflow.com',
    '$2a$12$yOr2E4TIj7DBIIhzPVnAOOWSy1GLhuVqqHcYJnZMzLtGfOh8rRbau',
    'Admin',
    'admin',
    TRUE
) ON CONFLICT (email) DO NOTHING;
