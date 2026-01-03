ALTER TABLE users ALTER COLUMN profile_picture DROP DEFAULT;

ALTER TABLE users
ALTER COLUMN profile_picture TYPE JSONB
USING CASE
    WHEN profile_picture IS NULL OR profile_picture = '' THEN '{"url": "", "public_id": ""}'::jsonb
    ELSE jsonb_build_object('url', profile_picture, 'public_id', '')
END;

ALTER TABLE users ALTER COLUMN profile_picture SET DEFAULT '{"url": "", "public_id": ""}'::jsonb;

ALTER TABLE jobs
ALTER COLUMN company_logo TYPE JSONB
USING CASE
    WHEN company_logo IS NULL OR company_logo = '' THEN '{"url": "", "public_id": ""}'::jsonb
    ELSE jsonb_build_object('url', company_logo, 'public_id', '')
END;
