DROP TABLE note_fts;

CREATE VIRTUAL TABLE note_fts USING fts4(id, markdown);
INSERT INTO note_fts SELECT id, markdown FROM note;