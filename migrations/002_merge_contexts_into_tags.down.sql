-- Recreate context tables (data loss is acceptable on rollback)
CREATE TABLE IF NOT EXISTS contexts (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    context_type TEXT NOT NULL DEFAULT 'location'
);

CREATE TABLE IF NOT EXISTS item_contexts (
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    context_id INT NOT NULL REFERENCES contexts(id) ON DELETE CASCADE,
    PRIMARY KEY (item_id, context_id)
);
