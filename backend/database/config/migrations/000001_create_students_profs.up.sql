BEGIN;

CREATE TABLE Students
(
    matric_no TEXT PRIMARY KEY,
    name      TEXT NOT NULL,
    email     TEXT
);

CREATE TABLE Professors
(
    id    BIGSERIAL PRIMARY KEY,
    name  TEXT NOT NULL,
    email TEXT
);

END;