ALTER TABLE users
ALTER COLUMN profile_picture TYPE TEXT
USING (profile_picture->>'url');

ALTER TABLE jobs
ALTER COLUMN company_logo TYPE VARCHAR(255)
USING (company_logo->>'url');
