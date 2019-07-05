
DO 
$$
BEGIN
-- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

SET CONSTRAINTS app_user_role_fk_u DEFERRED;
-- SET CONSTRAINTS app_user_role_fk_a DEFERRED;
SET CONSTRAINTS app_user_role_fk_ar DEFERRED;
SET CONSTRAINTS role_fk DEFERRED;

-- данные

INSERT INTO app (appname, description) VALUES ('app1'          , 'Тестовое express приложение 1');
INSERT INTO app (appname, description) VALUES ('app2'          , 'Тестовое express приложение 2');
INSERT INTO app (appname, description) VALUES ('onlinebc_admin', 'Тестовые трансляции'          );
INSERT INTO app (appname, description) VALUES ('rg'            , 'Сайт rg.ru'                   );

INSERT INTO "user" (username, password, email, fullname, description) VALUES ('vadim', '1', 'ivlev@rg.ru' , 'Ивлев Вадим'  , 'разработчик');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('serg' , '1', 'barsuk@rg.ru', 'Барсук Сергей', 'разработчик');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('max'  , '1', 'chagin@rg.ru', 'Чагин Максим' , 'начальник отдела разработки');

INSERT INTO app_role (appname, rolename, description) VALUES ('app1'          , 'manager','менеджер тестового приложения 1');
INSERT INTO app_role (appname, rolename, description) VALUES ('app1'          , 'worker' ,'работник тестового приложения 1');
INSERT INTO app_role (appname, rolename, description) VALUES ('app1'          , 'boss'   ,'босс тестового приложения 1');
INSERT INTO app_role (appname, rolename, description) VALUES ('app2'          , 'admin'  ,'админ тестового приложения 2');
INSERT INTO app_role (appname, rolename, description) VALUES ('app2'          , 'user'   ,'пользователь тестового приложения 2');
INSERT INTO app_role (appname, rolename, description) VALUES ('onlinebc_admin', 'admin'  ,'админ трансляций');
INSERT INTO app_role (appname, rolename, description) VALUES ('onlinebc_admin', 'editor' ,'редактор трансляций');
INSERT INTO app_role (appname, rolename, description) VALUES ('onlinebc_admin', 'guest'  ,'гость трансляций');
INSERT INTO app_role (appname, rolename, description) VALUES ('rg'            , 'reader' ,'читатель РГ');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('app1', 'vadim', 'manager');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('app1', 'vadim', 'worker');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('app1', 'vadim', 'boss');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('app1', 'max'  , 'manager');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('app1', 'max'  , 'boss');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('app1', 'serg' , 'manager');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('app2', 'vadim', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('app2', 'vadim', 'user');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('app2', 'max'  , 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('app2', 'serg' , 'user');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'vadim', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'vadim', 'editor');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'vadim', 'guest');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'max'  , 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'max'  , 'editor');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('onlinebc_admin', 'serg' , 'editor');


EXCEPTION WHEN OTHERS THEN 
    RAISE EXCEPTION 'Тестовые данные уже существуют.';
END;
$$;
