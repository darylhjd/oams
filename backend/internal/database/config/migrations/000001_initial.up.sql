BEGIN;

CREATE TYPE CLASS_TYPE AS ENUM ('LEC', 'TUT', 'LAB');

CREATE TYPE USER_ROLE AS ENUM ('STUDENT', 'COURSE_COORDINATOR', 'ADMIN');

CREATE TABLE users
(
    id         TEXT PRIMARY KEY, -- VCS Account No.
    name       TEXT      NOT NULL,
    email      TEXT      NOT NULL,
    role       USER_ROLE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE classes
(
    id         BIGSERIAL PRIMARY KEY,
    code       TEXT      NOT NULL,
    year       INTEGER   NOT NULL,
    semester   TEXT      NOT NULL,
    programme  TEXT      NOT NULL,
    au         SMALLINT  NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT ux_code_year_semester
        UNIQUE (code, year, semester)
);

CREATE TABLE class_groups
(
    id         BIGSERIAL PRIMARY KEY,
    class_id   BIGSERIAL  NOT NULL,
    name       TEXT       NOT NULL,
    class_type CLASS_TYPE NOT NULL,
    created_at TIMESTAMP  NOT NULL,
    updated_at TIMESTAMP  NOT NULL,
    CONSTRAINT ux_class_id_name
        UNIQUE (class_id, name),
    CONSTRAINT fk_class_id
        FOREIGN KEY (class_id)
            REFERENCES classes (id)
);

CREATE TABLE class_group_sessions
(
    id             BIGSERIAL PRIMARY KEY,
    class_group_id BIGSERIAL NOT NULL,
    start_time     TIMESTAMP NOT NULL,
    end_time       TIMESTAMP NOT NULL,
    venue          TEXT      NOT NULL,
    created_at     TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP NOT NULL,
    CONSTRAINT ux_class_group_id_start_time
        UNIQUE (class_group_id, start_time),
    CONSTRAINT fk_class_group_id
        FOREIGN KEY (class_group_id)
            REFERENCES class_groups (id)
);

CREATE TABLE session_enrollments
(
    session_id BIGSERIAL NOT NULL,
    user_id    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (session_id, user_id),
    CONSTRAINT fk_session_id
        FOREIGN KEY (session_id)
            REFERENCES class_group_sessions (id),
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);

COMMIT;