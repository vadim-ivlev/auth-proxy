
-- DO 
-- $$
-- BEGIN;
-- -- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

-- SET CONSTRAINTS app_user_role_fk_u DEFERRED;
-- SET CONSTRAINTS app_user_role_fk_a DEFERRED;

-- данные


INSERT INTO app (appname, url, description)                 VALUES ('auth'       ,'', 'Сервис авторизации');
INSERT INTO app (appname, url, description)                 VALUES ('node1'      ,'http://node:3001', 'Работает на продакшн env=prod. Служит для показа заголовков и тела запросов пропущенных через auth-proxy. Express приложение. test0');
INSERT INTO app (appname, url, description)                 VALUES ('node2'      ,'http://localhost:3002', 'Работает на компьютере разработчика env=dev. Служит для показа заголовков и тела запросов пропущенных через auth-proxy. Express приложение. test0');
INSERT INTO app (appname, url, description, rebase, public) VALUES ('pravo_rg_ru','https://pravo.rg.ru', 'Прокси к https://pravo.rg.ru . test0' ,'Y', 'Y' );
INSERT INTO app (appname, url, description, public)         VALUES ('rg'         ,'https://rg.ru'      , 'Прокси к https://rg.ru . test0', 'Y' );
INSERT INTO app (appname, url, description, public)         VALUES ('photoreports-admin'  ,'http://172.20.0.1:8091', 'GraphQL API админки фоторепов. test0', 'Y');


INSERT INTO "user" (username, password, email, fullname, description) VALUES ('admin', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'admin@rg.ru' , 'Админ Админов'  , 'Администратор auth-proxy. test0');
-- rosgas2011 => '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f'
-- 1 => '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b'


INSERT INTO app_user_role (appname, username, rolename) VALUES ('auth', 'admin', 'authadmin');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'admin', 'manager');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'admin', 'worker');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'admin', 'boss');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('node2', 'admin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node2', 'admin', 'user');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('pravo_rg_ru', 'admin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('pravo_rg_ru', 'admin', 'editor');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('pravo_rg_ru', 'admin', 'guest');


-- EXCEPTION WHEN OTHERS THEN 
--     RAISE EXCEPTION 'Тестовые данные уже существуют.';
-- END;
-- $$;
-- COMMIT;
