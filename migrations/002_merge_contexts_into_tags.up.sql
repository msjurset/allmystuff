-- Merge existing context names into tags
INSERT INTO tags (name)
SELECT DISTINCT c.name FROM contexts c
ON CONFLICT (name) DO NOTHING;

-- Migrate item_contexts relationships to item_tags
INSERT INTO item_tags (item_id, tag_id)
SELECT ic.item_id, t.id
FROM item_contexts ic
JOIN contexts c ON ic.context_id = c.id
JOIN tags t ON t.name = c.name
ON CONFLICT DO NOTHING;

-- Drop context tables
DROP TABLE item_contexts;
DROP TABLE contexts;
