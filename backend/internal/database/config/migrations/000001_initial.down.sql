BEGIN;

DROP TABLE session_enrollments;

DROP TABLE class_group_sessions;

DROP TABLE class_group_managers;

DROP TABLE class_groups;

DROP TABLE classes;

DROP TABLE user_signatures;

DROP TABLE users;

DROP FUNCTION update_updated_at;

DROP TYPE MANAGING_ROLE;

DROP TYPE USER_ROLE;

DROP TYPE CLASS_TYPE;

COMMIT;