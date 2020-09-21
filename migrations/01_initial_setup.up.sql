CREATE TABLE note (
    id INTEGER PRIMARY KEY,
    -- RFC3339 formatted
    date_created TEXT NOT NULL,
    date_updated TEXT,
    markdown TEXT NOT NULL,
    html TEXT NOT NULL
);

CREATE VIRTUAL TABLE note_fts USING fts4(id, markdown);
INSERT INTO note_fts SELECT id, markdown FROM note;
-- SELECT * FROM note_fts WHERE markdown MATCH 'link';

CREATE TRIGGER note_fts_before_delete BEFORE DELETE ON note
BEGIN
    DELETE FROM note_fts WHERE note_fts.id = OLD.id;
END;

CREATE TRIGGER note_fts_after_insert AFTER INSERT ON note
BEGIN
    INSERT INTO note_fts(id, markdown) VALUES(NEW.id, NEW.markdown);
END;

CREATE TRIGGER note_fts_after_update AFTER UPDATE ON note
BEGIN
    UPDATE note_fts SET markdown = NEW.markdown WHERE note_fts.id = NEW.id;
END;
