ALTER TABLE test_migrations ADD COLUMN i int;
---- create above / drop below ----
ALTER TABLE test_migrations DROP COLUMN i;
