BEGIN;

CREATE TYPE CLASS_TYPE AS ENUM ('LEC', 'TUT', 'LAB');

CREATE TYPE USER_ROLE AS ENUM ('USER', 'SYSTEM_ADMIN', 'EXTERNAL_SERVICE');

CREATE TYPE MANAGING_ROLE AS ENUM ('TEACHING_ASSISTANT', 'COURSE_COORDINATOR');

CREATE FUNCTION update_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    IF NEW <> OLD THEN
        NEW.updated_at := NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE users
(
    id         TEXT PRIMARY KEY, -- VCS Account No.
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL,
    role       USER_ROLE   NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE TABLE user_signatures
(
    user_id    TEXT        NOT NULL PRIMARY KEY,
    signature  TEXT        NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);

CREATE TRIGGER update_updated_at
    BEFORE UPDATE
    ON user_signatures
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE TABLE classes
(
    id         BIGSERIAL PRIMARY KEY,
    code       TEXT        NOT NULL,
    year       INTEGER     NOT NULL,
    semester   TEXT        NOT NULL,
    programme  TEXT        NOT NULL,
    au         SMALLINT    NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ux_code_year_semester
        UNIQUE (code, year, semester)
);

CREATE TRIGGER update_updated_at
    BEFORE UPDATE
    ON classes
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE TABLE class_attendance_rules
(
    id          BIGSERIAL PRIMARY KEY,
    class_id    BIGINT      NOT NULL,
    creator_id  TEXT        NOT NULL,
    title       TEXT        NOT NULL,
    description TEXT        NOT NULL,
    rule        TEXT        NOT NULL,
    environment JSONB       NOT NULL,
    active      BOOLEAN     NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ux_class_id_title
        UNIQUE (class_id, title),
    CONSTRAINT fk_class_id
        FOREIGN KEY (class_id)
            REFERENCES classes (id),
    CONSTRAINT fk_creator_id
        FOREIGN KEY (creator_id)
            REFERENCES users (id)
);

CREATE TRIGGER update_updated_at
    BEFORE UPDATE
    ON class_attendance_rules
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE TABLE class_groups
(
    id         BIGSERIAL PRIMARY KEY,
    class_id   BIGINT      NOT NULL,
    name       TEXT        NOT NULL,
    class_type CLASS_TYPE  NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ux_class_id_name_class_type
        UNIQUE (class_id, name, class_type),
    CONSTRAINT fk_class_id
        FOREIGN KEY (class_id)
            REFERENCES classes (id)
);

CREATE TRIGGER update_updated_at
    BEFORE UPDATE
    ON class_groups
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE TABLE class_group_managers
(
    id             BIGSERIAL PRIMARY KEY,
    user_id        TEXT          NOT NULL,
    class_group_id BIGINT        NOT NULL,
    managing_role  MANAGING_ROLE NOT NULL,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT ux_user_id_class_group_id
        UNIQUE (user_id, class_group_id),
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id),
    CONSTRAINT fk_class_group_id
        FOREIGN KEY (class_group_id)
            REFERENCES class_groups (id)
);

CREATE TRIGGER update_updated_at
    BEFORE UPDATE
    ON class_group_managers
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE TABLE class_group_sessions
(
    id             BIGSERIAL PRIMARY KEY,
    class_group_id BIGINT      NOT NULL,
    start_time     TIMESTAMPTZ NOT NULL,
    end_time       TIMESTAMPTZ NOT NULL,
    venue          TEXT        NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ux_class_group_id_start_time
        UNIQUE (class_group_id, start_time),
    CONSTRAINT fk_class_group_id
        FOREIGN KEY (class_group_id)
            REFERENCES class_groups (id),
    CONSTRAINT ck_start_time_more_than_end_time
        CHECK (start_time < end_time)
);

CREATE TRIGGER update_updated_at
    BEFORE UPDATE
    ON class_group_sessions
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE TABLE session_enrollments
(
    id         BIGSERIAL PRIMARY KEY,
    session_id BIGINT      NOT NULL,
    user_id    TEXT        NOT NULL,
    attended   BOOLEAN     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ux_session_id_user_id
        UNIQUE (session_id, user_id),
    CONSTRAINT fk_session_id
        FOREIGN KEY (session_id)
            REFERENCES class_group_sessions (id),
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);

CREATE TRIGGER update_updated_at
    BEFORE UPDATE
    ON session_enrollments
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

COMMIT;