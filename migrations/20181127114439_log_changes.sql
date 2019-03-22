-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE SCHEMA IF NOT EXISTS logging;

CREATE TABLE IF NOT EXISTS logging.t_history (
    id          serial,
    tstamp      timestamp DEFAULT now(),
    schemaname  text,
    tabname     text,
    operation   text,
    who         text DEFAULT current_user,
    new_val     json,
    old_val     json
);

CREATE FUNCTION change_trigger() RETURNS trigger AS $$ BEGIN IF TG_OP = 'INSERT' THEN INSERT INTO logging.t_history (tabname, schemaname, operation, new_val) VALUES (TG_RELNAME, TG_TABLE_SCHEMA, TG_OP, row_to_json(NEW));RETURN NEW;ELSIF   TG_OP = 'UPDATE'THEN INSERT INTO logging.t_history (tabname, schemaname, operation, new_val, old_val) VALUES (TG_RELNAME, TG_TABLE_SCHEMA, TG_OP, row_to_json(NEW), row_to_json(OLD)); RETURN NEW; ELSIF   TG_OP = 'DELETE' THEN INSERT INTO logging.t_history (tabname, schemaname, operation, old_val) VALUES (TG_RELNAME, TG_TABLE_SCHEMA, TG_OP, row_to_json(OLD)); RETURN OLD; END IF; END; $$ LANGUAGE 'plpgsql' SECURITY DEFINER;

-- Add log to one table
DROP TRIGGER IF EXISTS scooter_status_changes ON stooter_status;
CREATE TRIGGER scooter_status_changes BEFORE INSERT OR UPDATE OR DELETE ON scooter_status
    FOR EACH ROW EXECUTE PROCEDURE change_trigger();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TRIGGER IF EXISTS scooter_status_changes ON stooter_status;
DROP FUNCTION IF EXISTS change_trigger();
DROP TABLE IF EXISTS logging.t_history;
DROP SCHEMA IF EXISTS logging;
