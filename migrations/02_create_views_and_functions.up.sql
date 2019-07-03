
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



