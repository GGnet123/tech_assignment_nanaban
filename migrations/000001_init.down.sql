DROP FUNCTION IF EXISTS app.attach_updated_at_trigger(TEXT);

DROP FUNCTION IF EXISTS app.set_updated_at();

DROP SCHEMA IF EXISTS app CASCADE;

DO $$ BEGIN
  EXECUTE format('ALTER DATABASE %I SET search_path TO public', current_database());
END $$;

RESET TIME ZONE;