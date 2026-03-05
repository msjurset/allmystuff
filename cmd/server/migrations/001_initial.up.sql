CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    brand TEXT NOT NULL DEFAULT '',
    model TEXT NOT NULL DEFAULT '',
    serial_number TEXT NOT NULL DEFAULT '',
    purchase_date DATE,
    purchase_price NUMERIC(10,2),
    estimated_value NUMERIC(10,2),
    condition TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS item_tags (
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (item_id, tag_id)
);

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

CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    filepath TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_images_item_id ON images(item_id);
