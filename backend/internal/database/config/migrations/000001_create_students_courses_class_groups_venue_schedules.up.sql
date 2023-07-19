BEGIN;

CREATE TABLE Students
(
    id    TEXT PRIMARY KEY,
    name  TEXT NOT NULL,
    email TEXT
);

CREATE TABLE Courses
(
    code      TEXT PRIMARY KEY,
    programme TEXT     NOT NULL,
    au        SMALLINT NOT NULL
);

CREATE TABLE Class_Groups
(
    id          BIGSERIAL PRIMARY KEY,
    course_code TEXT    NOT NULL,
    year        INTEGER NOT NULL,
    semester    TEXT    NOT NULL,
    group_name  TEXT    NOT NULL,
    UNIQUE (course_code, year, semester, group_name),
    CONSTRAINT fk_course_code
        FOREIGN KEY (course_code)
            REFERENCES Courses (code)
);

CREATE TABLE Class_Venue_Schedules
(
    class_group_id BIGSERIAL,
    datetime       TIMESTAMP,
    venue          TEXT NOT NULL,
    PRIMARY KEY (class_group_id, datetime),
    CONSTRAINT fk_class_group_id
        FOREIGN KEY (class_group_id)
            REFERENCES Class_Groups (id)
);

END;