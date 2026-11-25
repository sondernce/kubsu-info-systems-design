CREATE TABLE drivers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    last_name TEXT NOT NULL,
    first_name TEXT NOT NULL,
    middle_name TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    experience INTEGER NOT NULL CHECK (experience BETWEEN 0 AND 70)
);
