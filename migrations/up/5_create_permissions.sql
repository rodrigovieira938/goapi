INSERT INTO permission (name, description)
VALUES
    ('write:cars', 'Create, update, or delete car records'),
    ('read:users', 'View user information and lists'),
    ('write:users', 'Grant or revoke user permissions'),
    ('read:reservations', 'View all reservations (not only own)'),
    ('write:reservations', 'Create, update, or delete reservations')
ON CONFLICT DO NOTHING;
