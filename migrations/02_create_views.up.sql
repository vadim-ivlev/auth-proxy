
-- V I E W S ------------------------------------------------------------------

-- Расширенное представление основной таблицы
-- С дополнительными полями из справочных таблиц 
DROP VIEW IF EXISTS app_user_role_extended;
CREATE VIEW app_user_role_extended AS
    SELECT 
        aur.appname    AS appname,
        aur.username   AS username,
        aur.rolename   AS rolename,
        u.email        AS user_email,
        u.fullname     AS user_fullname,
        u.description  AS user_description,
        u.disabled     AS user_disabled,
        a.description  AS app_description,
        a.url          AS app_url
        
    
    FROM app_user_role   AS aur
    INNER JOIN "user"    AS u   ON aur.username = u.username
    INNER JOIN app       AS a   ON aur.appname  = a.appname
;

DROP VIEW IF EXISTS group_user_role_extended;
CREATE VIEW group_user_role_extended AS
    SELECT 
        gu.group_id   AS group_id,
        gu.user_id    AS user_id,
        gu.rolename   AS rolename,

        g.groupname   AS group_groupname,
        g.description AS group_description,

        u.email       AS user_email,
        u.fullname    AS user_fullname,
        u.description AS user_description,
        u.disabled    AS user_disabled
        
    FROM "group" g
    JOIN group_user_role gu ON gu.group_id = g.id
    JOIN "user" u           ON u.id = gu.user_id
;

DROP VIEW IF EXISTS group_app_role_extended;
CREATE VIEW group_app_role_extended AS
    SELECT 
        ga.group_id   AS group_id,
        ga.app_id     AS app_id,
        ga.rolename   AS rolename,

        g.groupname   AS group_groupname,
        g.description AS group_description,

        a.description  AS app_description,
        a.url          AS app_url
           
    FROM "group" g
    JOIN group_app_role ga ON ga.group_id = g.id
    JOIN app a             ON a.id = ga.app_id
;
