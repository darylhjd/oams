BEGIN;

DROP TABLE session_enrollments;

DROP TABLE class_group_sessions;

DROP TABLE class_groups;

DROP TABLE class_managers;

DROP TABLE classes;

DROP TABLE users;

DROP FUNCTION update_updated_at;

DROP TYPE MANAGING_ROLE;

DROP TYPE USER_ROLE;

DROP TYPE CLASS_TYPE;

COMMIT;