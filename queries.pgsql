-- SELECT * FROM full_user;


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
            ) as subtable ;



SELECT jsonb_agg(subtable) from  
            ( 
            SELECT DISTINCT 
                aur.appname,
                (
                SELECT 
                    jsonb_agg(roles_subtable) 
                FROM  
                        ( 
                        SELECT a.rolename
                        FROM app_user_role AS a
                        WHERE a.username = aur.username 
                        AND   a.appname  = aur.appname
                        ORDER BY a.rolename ASC
                        ) as roles_subtable 
                ) as roles

            FROM app_user_role AS aur
            WHERE aur.username = 'vadim'
            ORDER BY aur.appname ASC
            ) as subtable ;





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
-- WHERE aur.username = 'vadim'
ORDER BY aur.username,aur.appname ASC;


