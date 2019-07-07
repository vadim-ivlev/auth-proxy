-- TODO: remove
-- CREATE OR REPLACE FUNCTION get_app_user_roles(a_name text, u_name text)
--  RETURNS json
--  LANGUAGE plpgsql
-- AS $function$
-- BEGIN
--     RETURN
--     (
--             SELECT jsonb_agg(rolename) as roles FROM app_user_role
--             WHERE  appname  = a_name 
--             AND    username = u_name
--     );

-- END;
-- $function$
-- ;


-- V I E W S ------------------------------------------------------------------

-- Расширенное представление основной таблицы
-- С дополнительными полями из справочных таблиц 
CREATE OR REPLACE VIEW app_user_role_extended AS
    SELECT 
        aur.appname    AS appname,
        aur.username   AS username,
        aur.rolename   AS rolename,
        u.email        AS user_email,
        u.fullname     AS user_fullname,
        u.description  AS user_description,
        a.description  AS app_description,
        ar.description AS app_role_description
    
    FROM app_user_role   AS aur
    INNER JOIN "user"    AS u   ON aur.username = u.username
    INNER JOIN app       AS a   ON aur.appname  = a.appname
    INNER JOIN app_role  AS ar  ON aur.appname = ar.appname AND aur.rolename = ar.rolename
;



CREATE OR REPLACE VIEW full_app AS
    SELECT  *,  
        ( 
            select jsonb_agg(subtable) from  
            ( SELECT * FROM app_role WHERE appname = app.appname ) 
            as subtable 
        ) AS roles
    FROM app 
;

-- Пользователь с выборкой его ролей в приложениях 
-- в виде JSON поля
CREATE OR REPLACE VIEW full_user AS
    SELECT  *,  
        ( SELECT json_agg(subtable) from  
            ( SELECT *  FROM app_user_role 
            WHERE username = "user".username 
            ) as subtable 
        ) AS app_user_roles
    FROM "user" 
;

-- Пользователь с выборкой его приложений и его ролей в них 
-- в виде JSON поля
CREATE OR REPLACE VIEW full_user2 AS
    SELECT  *,  
        ( 
        SELECT jsonb_agg(subtable) from  
            ( 
            SELECT DISTINCT 
                aur.appname,
                aur.app_description,
                (
                SELECT 
                    jsonb_agg(roles_subtable) 
                FROM  
                        ( 
                        SELECT 
                        a.rolename,
                        a.app_role_description
                        FROM app_user_role_extended AS a
                        WHERE a.username = aur.username 
                        AND   a.appname  = aur.appname
                        ORDER BY a.rolename ASC
                        ) as roles_subtable 
                ) as roles

            FROM app_user_role_extended AS aur
            WHERE aur.username = 'vadim'
            ORDER BY aur.appname ASC
            ) as subtable 
        ) AS apps
    FROM "user" 
;









-- TODO: delete
-- Расширенное представление основной таблицы
-- С дополнительными полями из справочных таблиц, 
-- И JSON полем roles c ролями пользователя 
-- для каждой пары (username, appname)
CREATE OR REPLACE VIEW user_roles AS
    SELECT DISTINCT
        aur.*,
        (
        SELECT 
            jsonb_agg(roles_subtable) 
        FROM  
                ( 
                SELECT 
                a.rolename,
                a.app_role_description
                FROM app_user_role_extended AS a
                WHERE a.username = aur.username 
                AND   a.appname  = aur.appname
                ORDER BY a.rolename ASC
                ) as roles_subtable 
        ) as roles

    FROM app_user_role_extended AS aur
    ORDER BY aur.username,aur.appname ASC
;
