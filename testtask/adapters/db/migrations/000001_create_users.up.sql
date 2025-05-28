CREATE TABLE IF NOT EXISTS people (
    people_id VARCHAR(255) PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT,
    age INTEGER NOT NULL CHECK (age > 0),
    gender TEXT NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    nationality TEXT NOT NULL
);