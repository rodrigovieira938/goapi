CREATE TABLE IF NOT EXISTS reservation (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL,
    car_id     INT NOT NULL,
    start_date DATE NOT NULL,
    end_date   DATE NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES "user"(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_car
        FOREIGN KEY(car_id)
        REFERENCES car(id)
        ON DELETE CASCADE
);
