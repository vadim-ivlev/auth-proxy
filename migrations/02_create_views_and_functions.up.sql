
CREATE OR REPLACE FUNCTION public.get_app_user_roles(a_name text, u_name text)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN
    (
        select jsonb_agg(t) from
        (
            SELECT rolename FROM app_user_role
            WHERE  appname  = a_name 
            AND    username = u_name
        ) t
    );

END;
$function$
;



