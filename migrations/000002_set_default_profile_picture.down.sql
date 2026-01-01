ALTER TABLE users ALTER COLUMN profile_picture DROP NOT NULL;
ALTER TABLE users ALTER COLUMN profile_picture DROP DEFAULT;
-- We cannot easily restore original NULLs unless we knew which ones were NULL before. 
-- But generally down migration reverts schema constraints.
