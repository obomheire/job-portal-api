UPDATE users SET profile_picture = '' WHERE profile_picture IS NULL;
ALTER TABLE users ALTER COLUMN profile_picture SET DEFAULT '';
ALTER TABLE users ALTER COLUMN profile_picture SET NOT NULL;
