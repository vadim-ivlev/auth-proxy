
DO 
$$
BEGIN
-- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

SET CONSTRAINTS app_user_role_fk_a DEFERRED;
SET CONSTRAINTS app_user_role_fk_u DEFERRED;
-- SET CONSTRAINTS app_user_role_fk_r DEFERRED;
SET CONSTRAINTS role_fk DEFERRED;

-- данные

INSERT INTO public.app (appname, description, url) VALUES ('app1'          , 'Тестовое express приложение 1', 'http://localhost:3001');
INSERT INTO public.app (appname, description, url) VALUES ('app2'          , 'Тестовое express приложение 2', 'http://localhost:3002');
INSERT INTO public.app (appname, description, url) VALUES ('onlinebc_admin', 'Тестовые трансляции'          , 'http://localhost:7700');
INSERT INTO public.app (appname, description, url) VALUES ('rg'            , 'Сайт rg.ru'                   , 'https://rg.ru');

INSERT INTO public.user (username, password, email, fullname, description) VALUES ('vadim', '1', 'ivlev@rg.ru' , 'Ивлев Вадим'  , 'разработчик');
INSERT INTO public.user (username, password, email, fullname, description) VALUES ('serg' , '1', 'barsuk@rg.ru', 'Барсук Сергей', 'разработчик');
INSERT INTO public.user (username, password, email, fullname, description) VALUES ('max'  , '1', 'chagin@rg.ru', 'Чагин Максим' , 'начальник отдела разработки');

INSERT INTO public.role (appname, rolename) VALUES ('app1'          , 'manager');
INSERT INTO public.role (appname, rolename) VALUES ('app1'          , 'worker');
INSERT INTO public.role (appname, rolename) VALUES ('app1'          , 'boss');
INSERT INTO public.role (appname, rolename) VALUES ('app2'          , 'admin');
INSERT INTO public.role (appname, rolename) VALUES ('app2'          , 'user');
INSERT INTO public.role (appname, rolename) VALUES ('onlinebc_admin', 'admin');
INSERT INTO public.role (appname, rolename) VALUES ('onlinebc_admin', 'editor');
INSERT INTO public.role (appname, rolename) VALUES ('onlinebc_admin', 'guest');
INSERT INTO public.role (appname, rolename) VALUES ('rg'            , 'reader');

INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('app1', 'vadim', 'manager');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('app1', 'vadim', 'worker');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('app1', 'vadim', 'boss');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('app1', 'max'  , 'manager');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('app1', 'max'  , 'boss');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('app1', 'serg' , 'manager');

INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('app2', 'vadim', 'admin');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('app2', 'vadim', 'user');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('app2', 'max'  , 'admin');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('app2', 'serg' , 'user');

INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'vadim', 'admin');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'vadim', 'editor');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'vadim', 'guest');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'max'  , 'admin');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'max'  , 'editor');
INSERT INTO public.app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'serg' , 'editor');


EXCEPTION WHEN OTHERS THEN 
    RAISE EXCEPTION 'Тестовые данные уже существуют.';
END;
$$;
