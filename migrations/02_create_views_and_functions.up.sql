-- TODO: remove
CREATE OR REPLACE FUNCTION public.get_app_user_roles(a_name text, u_name text)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN
    (
            SELECT jsonb_agg(rolename) as roles FROM app_user_role
            WHERE  appname  = a_name 
            AND    username = u_name
    );

END;
$function$
;

CREATE OR REPLACE VIEW public.full_app AS
    SELECT  *,  
        ( 
            select jsonb_agg(subtable) from  
            ( SELECT * FROM app_role WHERE appname = app.appname ) 
            as subtable 
        ) AS roles
    FROM app 
;

CREATE OR REPLACE VIEW public.full_user AS
SELECT  *,  
    ( SELECT json_agg(subtable) from  
        ( SELECT *  FROM app_user_role 
          WHERE username = "user".username 
        ) as subtable 
    ) AS app_user_roles
FROM "user" 
;
