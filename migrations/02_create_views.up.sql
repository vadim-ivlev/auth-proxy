
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

