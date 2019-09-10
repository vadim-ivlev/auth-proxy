
-- DO 
-- $$
-- BEGIN;
-- -- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

-- SET CONSTRAINTS app_user_role_fk_u DEFERRED;
-- SET CONSTRAINTS app_user_role_fk_a DEFERRED;

-- данные


INSERT INTO app (appname, url, description)         VALUES ('auth'       ,'', 'Сервис авторизации');
INSERT INTO app (appname, url, description)         VALUES ('node1'      ,'http://node:3001', 'Работает на продакшн env=prod. Служит для показа заголовков и тела запросов пропущенных через auth-proxy. Express приложение. test0');
INSERT INTO app (appname, url, description)         VALUES ('node2'      ,'http://localhost:3002', 'Работает на компьютере разработчика env=dev. Служит для показа заголовков и тела запросов пропущенных через auth-proxy. Express приложение. test0');
INSERT INTO app (appname, url, description, rebase) VALUES ('pravo_rg_ru','https://pravo.rg.ru', 'Прокси к https://pravo.rg.ru . test0' ,'Y' );
INSERT INTO app (appname, url, description)         VALUES ('rg'         ,'https://rg.ru'      , 'Прокси к https://rg.ru . test0' );
INSERT INTO app (appname, url, description)         VALUES ('wikipedia'  ,'https://www.wikipedia.org', 'Прокси к Википедии. test0');


INSERT INTO "user" (username, password, email, fullname, description) VALUES ('vadim', '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b', 'ivlev@rg.ru' , 'Ивлев Вадим'  , 'разработчик test0');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('max'  , '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b', 'chagin@rg.ru', 'Чагин Максим' , 'начальник отдела разработки test0');


INSERT INTO app_user_role (appname, username, rolename) VALUES ('auth', 'vadim', 'authadmin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('auth', 'max', 'authadmin');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'vadim', 'manager');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'vadim', 'worker');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'vadim', 'boss');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'max'  , 'manager');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'max'  , 'boss');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('node2', 'vadim', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node2', 'vadim', 'user');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node2', 'max'  , 'admin');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('pravo_rg_ru', 'vadim', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('pravo_rg_ru', 'vadim', 'editor');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('pravo_rg_ru', 'vadim', 'guest');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('pravo_rg_ru', 'max'  , 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('pravo_rg_ru', 'max'  , 'editor');


-- EXCEPTION WHEN OTHERS THEN 
--     RAISE EXCEPTION 'Тестовые данные уже существуют.';
-- END;
-- $$;
-- COMMIT;
