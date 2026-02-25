-- B1 Identity: verified profile. email_verified already exists; add phone_verified for badge.
ALTER TABLE users ADD COLUMN phone_verified INTEGER NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN phone TEXT;
