
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS whitelists
(
  id serial NOT NULL,
  username character varying(255) NOT NULL,
  createdat timestamp without time zone DEFAULT now(),
  profilesid integer,
  CONSTRAINT whitelists_pkey PRIMARY KEY (id ),
  CONSTRAINT whitelists_profilesid_fkey FOREIGN KEY (profilesid)
      REFERENCES profiles (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION notify_trigger_whitelists()
  RETURNS trigger AS
$BODY$BEGIN IF (TG_OP = 'DELETE') THEN PERFORM pg_notify('notify_trigger_whitelists', '[{' || TG_TABLE_NAME || ':' || OLD.id || '}, { operation: "' || TG_OP || '"}]');RETURN old;ELSIF (TG_OP = 'INSERT') THEN PERFORM pg_notify('notify_trigger_whitelists', '[{' || TG_TABLE_NAME || ':' || NEW.id || '}, { operation: "' || TG_OP || '"}]');RETURN new; ELSIF (TG_OP = 'UPDATE') THEN PERFORM pg_notify('notify_trigger_whitelists', '[{' || TG_TABLE_NAME || ':' || NEW.id || '}, { operation: "' || TG_OP || '"}]');RETURN new; END IF; END; $BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
-- +goose StatementEnd

CREATE TRIGGER watched_table_trigger
  AFTER INSERT OR UPDATE OR DELETE
  ON whitelists
  FOR EACH ROW
  EXECUTE PROCEDURE notify_trigger_whitelists();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE whitelists;
