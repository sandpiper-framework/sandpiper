CREATE DATABASE sandpiper;
DO $$
    BEGIN
        CREATE USER sandpiper WITH ENCRYPTED PASSWORD 'autocare';
    EXCEPTION WHEN DUPLICATE_OBJECT THEN
        RAISE NOTICE 'CREATE USER skipped: "sandpiper" user already exists.';
    END
$$;
GRANT ALL PRIVILEGES ON DATABASE sandpiper TO sandpiper;
