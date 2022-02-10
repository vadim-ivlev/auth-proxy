DROP INDEX IF EXISTS user_pashash_idx;
DROP INDEX IF EXISTS aur_appname_idx;
DROP INDEX IF EXISTS aur_username_idx;
DROP INDEX IF EXISTS user_email_unique_idx;
DROP INDEX IF EXISTS user_id_unique_idx;
-- DROP INDEX IF EXISTS app_textsearch_idx;
-- DROP INDEX IF EXISTS user_textsearch_idx;

DROP TABLE IF EXISTS app_user_role;
DROP TABLE IF EXISTS app;
DROP TABLE IF EXISTS "user";

DROP TABLE IF EXISTS group_user_role;
DROP TABLE IF EXISTS group_app_role;
DROP TABLE IF EXISTS "group";

DROP SCHEMA IF EXISTS auth CASCADE ;