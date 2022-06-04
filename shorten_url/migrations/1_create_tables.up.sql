CREATE TABLE url (
    id TEXT PRIMARY KEY,
    original_url TEXT NOT NULL,
    owner_by TEXT NULL DEFAULT 'guest'
);