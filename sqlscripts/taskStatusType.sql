DO
$$
BEGIN
  IF NOT EXISTS (SELECT *
                        FROM pg_type typ
                             INNER JOIN pg_namespace nsp
                                        ON nsp.oid = typ.typnamespace
                        WHERE nsp.nspname = current_schema()
                              AND typ.typname = 'task_status_type') THEN
    CREATE TYPE task_status_type AS ENUM ('TODO', 'DOING', 'DONE');
  END IF;
END;
$$
LANGUAGE plpgsql;