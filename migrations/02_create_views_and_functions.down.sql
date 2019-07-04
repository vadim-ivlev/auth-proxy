DROP VIEW IF EXISTS public.full_app;
DROP VIEW IF EXISTS public.full_user;

DROP FUNCTION IF EXISTS public.get_app_user_roles(a_name text, u_name text);
