INSERT INTO "user" (username, email, password)
VALUES ('admin','admin', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Grant all permissions to the admin user
DO $$
DECLARE
    admin_id INT;
    perm RECORD;
BEGIN
    -- Fetch admin user id
    SELECT id INTO admin_id FROM "user" WHERE username = 'admin';

    IF admin_id IS NOT NULL THEN
        -- Loop over all permissions and assign each to admin
        FOR perm IN SELECT id FROM permission LOOP
            INSERT INTO user_permission (user_id, permission_id)
            VALUES (admin_id, perm.id)
            ON CONFLICT DO NOTHING;
        END LOOP;
    END IF;
END$$;