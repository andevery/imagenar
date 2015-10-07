
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS tasks
(
  id serial NOT NULL,
  type integer DEFAULT 1,
  status integer DEFAULT 0,
  follows boolean DEFAULT false,
  likes boolean DEFAULT false,
  maxlikes integer DEFAULT 2,
  minlikes integer DEFAULT 4,
  maxtags integer DEFAULT 50,
  maxfollowedby integer DEFAULT 500,
  minfollowedby integer DEFAULT 0,
  maxfollows integer DEFAULT 300,
  minfollows integer DEFAULT 100,
  minmedia integer DEFAULT 20,
  delay integer DEFAULT 60,
  tags character varying(255),
  likescount integer DEFAULT 0,
  followscount integer DEFAULT 0,
  unfollowscount integer DEFAULT 0,
  createdat timestamp without time zone DEFAULT now(),
  profilesid integer,
  CONSTRAINT tasks_pkey PRIMARY KEY (id ),
  CONSTRAINT tasks_profilesid_fkey FOREIGN KEY (profilesid)
      REFERENCES profiles (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION tasks_notify_event() RETURNS TRIGGER AS $$

  DECLARE
      data json;
      notification json;
  BEGIN

      -- Convert the old or new row to JSON, based on the kind of action.
      -- Action = DELETE?             -> OLD row
      -- Action = INSERT or UPDATE?   -> NEW row
      IF (TG_OP = 'DELETE') THEN
          data = row_to_json(OLD);
      ELSE
          data = row_to_json(NEW);
      END IF;

      -- Contruct the notification as a JSON string.
      notification = json_build_object(
                        'table',TG_TABLE_NAME,
                        'action', TG_OP,
                        'data', data);


      -- Execute pg_notify(channel, notification)
      PERFORM pg_notify('tasks_notify_event',notification::text);

      -- Result is ignored since this is an AFTER trigger
      RETURN NULL;
  END;

$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION notify_trigger_tasks()
  RETURNS trigger AS
$BODY$BEGIN IF (TG_OP = 'DELETE') THEN PERFORM pg_notify('notify_trigger_tasks', '[{' || TG_TABLE_NAME || ':' || OLD.id || '}, { operation: "' || TG_OP || '"}]');RETURN old;ELSIF (TG_OP = 'INSERT') THEN PERFORM pg_notify('notify_trigger_tasks', '[{' || TG_TABLE_NAME || ':' || NEW.id || '}, { operation: "' || TG_OP || '"}]');RETURN new; ELSIF (TG_OP = 'UPDATE') THEN PERFORM pg_notify('notify_trigger_tasks', '[{' || TG_TABLE_NAME || ':' || NEW.id || '}, { operation: "' || TG_OP || '"}]');RETURN new; END IF; END; $BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
-- +goose StatementEnd

CREATE TRIGGER tasks_notify_event
AFTER INSERT OR UPDATE OR DELETE ON tasks
    FOR EACH ROW EXECUTE PROCEDURE tasks_notify_event();

CREATE TRIGGER watched_table_trigger
  AFTER INSERT OR UPDATE OR DELETE
  ON tasks
  FOR EACH ROW
  EXECUTE PROCEDURE notify_trigger_tasks();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE tasks;
