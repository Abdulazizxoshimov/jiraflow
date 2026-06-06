ALTER TABLE issues
    ADD COLUMN resolution VARCHAR(50)
        CHECK (resolution IN ('fixed', 'wont_fix', 'duplicate', 'incomplete', 'cannot_reproduce', 'done'));
