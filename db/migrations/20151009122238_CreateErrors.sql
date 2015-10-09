
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS errors
(
  id serial NOT NULL,
  type character varying(255),
  message character varying(255),
  createdat timestamp without time zone DEFAULT now(),
  tasksid integer,
  CONSTRAINT errors_pkey PRIMARY KEY (id ),
  CONSTRAINT errors_tasksid_fkey FOREIGN KEY (tasksid)
      REFERENCES tasks (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION notify_trigger_errors()
  RETURNS trigger AS
$BODY$BEGIN IF (TG_OP = 'DELETE') THEN PERFORM pg_notify('notify_trigger_errors', '[{' || TG_TABLE_NAME || ':' || OLD.id || '}, { operation: "' || TG_OP || '"}]');RETURN old;ELSIF (TG_OP = 'INSERT') THEN PERFORM pg_notify('notify_trigger_errors', '[{' || TG_TABLE_NAME || ':' || NEW.id || '}, { operation: "' || TG_OP || '"}]');RETURN new; ELSIF (TG_OP = 'UPDATE') THEN PERFORM pg_notify('notify_trigger_errors', '[{' || TG_TABLE_NAME || ':' || NEW.id || '}, { operation: "' || TG_OP || '"}]');RETURN new; END IF; END; $BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
-- +goose StatementEnd

CREATE TRIGGER watched_table_trigger
  AFTER INSERT OR UPDATE OR DELETE
  ON errors
  FOR EACH ROW
  EXECUTE PROCEDURE notify_trigger_errors();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE errors;
