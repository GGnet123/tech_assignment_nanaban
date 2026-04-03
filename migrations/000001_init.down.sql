DROP FUNCTION IF EXISTS attach_updated_at_trigger(TEXT);

DROP FUNCTION IF EXISTS set_updated_at();

DROP SCHEMA IF EXISTS app CASCADE;

SET search_path TO public;

RESET TIME ZONE;