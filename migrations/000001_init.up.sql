-- separate schema
CREATE SCHEMA IF NOT EXISTS app;

-- set default search path for the database permanently (future connections)
DO $$ BEGIN
  EXECUTE format('ALTER DATABASE %I SET search_path TO app', current_database());
END $$;

-- set search path for the current session
SET search_path TO app;

-- default timezone
SET TIME ZONE 'UTC';

-- func to auto set updated at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION attach_updated_at_trigger(table_name TEXT)
RETURNS void AS $$
BEGIN
EXECUTE format(
        'CREATE TRIGGER trg_%I_updated_at
         BEFORE UPDATE ON %I
         FOR EACH ROW
         EXECUTE FUNCTION set_updated_at()',
        table_name, table_name
        );
END;
$$ LANGUAGE plpgsql;