-- SQLITE VERSION *********************************************************************




DROP VIEW IF EXISTS full_app;
CREATE VIEW full_app AS
    SELECT  *,  
        ( select json_group_array(rolename) from  
            ( SELECT DISTINCT rolename FROM app_user_role WHERE appname = app.appname ) 
        ) AS roles
    FROM app 
;

-- Пользователь с выборкой его ролей в приложениях 
-- в виде JSON поля
DROP VIEW IF EXISTS full_user;
CREATE VIEW full_user AS
    SELECT  *,  
        ( select json_group_array(appname) from  
            ( SELECT DISTINCT appname  FROM app_user_role WHERE username = "user".username ) 
        ) AS apps
    FROM "user" 
;



