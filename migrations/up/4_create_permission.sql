-- Permissions table
CREATE TABLE IF NOT EXISTS permission (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT
);

-- User-Permission join table
CREATE TABLE IF NOT EXISTS user_permission (
    user_id       INT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (user_id, permission_id),
    CONSTRAINT fk_user_perm
        FOREIGN KEY(user_id)
        REFERENCES "user"(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_permission
        FOREIGN KEY(permission_id)
        REFERENCES permission(id)
        ON DELETE CASCADE
);
