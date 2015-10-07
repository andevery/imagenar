
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS profiles
(
  id serial NOT NULL,
  username character varying(255),
  password character varying(255),
  createdat timestamp without time zone DEFAULT now(),
  CONSTRAINT profiles_pkey PRIMARY KEY (id ),
  CONSTRAINT profiles_username_key UNIQUE (username )
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION notify_trigger_profiles()
  RETURNS trigger AS
$BODY$BEGIN IF (TG_OP = 'DELETE') THEN PERFORM pg_notify('notify_trigger_profiles', '[{' || TG_TABLE_NAME || ':' || OLD.id || '}, { operation: "' || TG_OP || '"}]');RETURN old;ELSIF (TG_OP = 'INSERT') THEN PERFORM pg_notify('notify_trigger_profiles', '[{' || TG_TABLE_NAME || ':' || NEW.id || '}, { operation: "' || TG_OP || '"}]');RETURN new; ELSIF (TG_OP = 'UPDATE') THEN PERFORM pg_notify('notify_trigger_profiles', '[{' || TG_TABLE_NAME || ':' || NEW.id || '}, { operation: "' || TG_OP || '"}]');RETURN new; END IF; END; $BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
-- +goose StatementEnd

CREATE TRIGGER watched_table_trigger
  AFTER INSERT OR UPDATE OR DELETE
  ON profiles
  FOR EACH ROW
  EXECUTE PROCEDURE notify_trigger_profiles();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE profiles;
