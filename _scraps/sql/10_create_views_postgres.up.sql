

DROP VIEW IF EXISTS full_app;
CREATE VIEW full_app AS
    SELECT  *,  
        ( select jsonb_agg(subtable) from  
            ( SELECT DISTINCT rolename FROM app_user_role WHERE appname = app.appname ) 
            as subtable 
        ) AS roles
    FROM app 
;

-- Пользователь с выборкой его ролей в приложениях 
-- в виде JSON поля
DROP VIEW IF EXISTS full_user;
CREATE VIEW full_user AS
    SELECT  *,  
        ( select json_agg(subtable) from  
            ( SELECT DISTINCT appname  FROM app_user_role WHERE username = "user".username ) 
            as subtable 
        ) AS apps
    FROM "user" 
;



-- -- TODO: delete
-- -- Пользователь с выборкой его приложений и его ролей в них 
-- -- в виде JSON поля
-- CREATE OR REPLACE VIEW full_user2 AS
--     SELECT  *,  
--         ( 
--         SELECT jsonb_agg(subtable) from  
--             ( 
--             SELECT DISTINCT 
--                 aur.appname,
--                 aur.app_description,
--                 (
--                 SELECT 
--                     jsonb_agg(roles_subtable) 
--                 FROM  
--                         ( 
--                         SELECT 
--                         a.rolename,
--                         -- a.app_role_description
--                         FROM app_user_role_extended AS a
--                         WHERE a.username = aur.username 
--                         AND   a.appname  = aur.appname
--                         ORDER BY a.rolename ASC
--                         ) as roles_subtable 
--                 ) as roles

--             FROM app_user_role_extended AS aur
--             WHERE aur.username = 'vadim'
--             ORDER BY aur.appname ASC
--             ) as subtable 
--         ) AS apps
--     FROM "user" 
-- ;



-- -- TODO: delete
-- -- Расширенное представление основной таблицы
-- -- С дополнительными полями из справочных таблиц, 
-- -- И JSON полем roles c ролями пользователя 
-- -- для каждой пары (username, appname)
-- CREATE OR REPLACE VIEW user_roles AS
--     SELECT DISTINCT
--         aur.*,
--         (
--         SELECT 
--             jsonb_agg(roles_subtable) 
--         FROM  
--                 ( 
--                 SELECT 
--                 a.rolename,
--                 -- a.app_role_description
--                 FROM app_user_role_extended AS a
--                 WHERE a.username = aur.username 
--                 AND   a.appname  = aur.appname
--                 ORDER BY a.rolename ASC
--                 ) as roles_subtable 
--         ) as roles

--     FROM app_user_role_extended AS aur
--     ORDER BY aur.username,aur.appname ASC
-- ;