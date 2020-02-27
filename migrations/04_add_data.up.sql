
-- DO 
-- $$
-- BEGIN;
-- -- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

-- SET CONSTRAINTS app_user_role_fk_u DEFERRED;
-- SET CONSTRAINTS app_user_role_fk_a DEFERRED;

-- данные


INSERT INTO app (appname, url, description) VALUES ('auth','', 'Сервис авторизации');
INSERT INTO app (appname, url, description, sign) VALUES ('node1','http://auth-node:3001', 'Работает на продакшн env=prod. Служит для показа заголовков и тела запросов пропущенных через auth-proxy. Express приложение', 'Y');
INSERT INTO app (appname, url, description, sign) VALUES ('node2','http://localhost:3002', 'Работает на компьютере разработчика env=dev. Служит для показа заголовков и тела запросов пропущенных через auth-proxy. Express приложение', 'Y');
INSERT INTO app (appname, url, description, public) VALUES ('rg','https://rg.ru', 'Прокси к https://rg.ru', 'Y' );
INSERT INTO app (appname, url, description, public) VALUES ('photoreports-admin','http://photoreports-admin:8091', 'GraphQL API админки фоторепов');
INSERT INTO app (appname, url, description, public) VALUES ('rgcore','http://rgru-core:8076', 'GraphQL API редакторского интерфейса');


INSERT INTO "user" (username, password, email, fullname, description) VALUES ('admin', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'admin@rg.ru' , 'Админ Админов'  , 'Администратор auth-proxy');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('test', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'dev@rg.ru' , 'Иван Иванов'  , 'Для тестирования auth-proxy');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('chagin', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'chagin@rg.ru' , 'Максим Чагин'  , 'Администратор auth-proxy');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('chernyshev', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'chernyshev@rg.ru' , 'Алексей Чернышев'  , 'Администратор auth-proxy');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ismagilov', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'ismagilov@rg.ru' , 'Камиль Исмагилов'  , 'Администратор auth-proxy');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kondratiev', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'kondratiev@rg.ru' , 'Антон Кондратьев'  , 'Администратор auth-proxy');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ivlev', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'ivlev@rg.ru' , 'Вадим Ивлев'  , 'Администратор auth-proxy');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('barsuk', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'barsuk@rg.ru' , 'Сергей Барсук'  , 'Администратор auth-proxy');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('filatchev', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'filatchev@rg.ru' , 'Дмитрий Филатчев'  , 'Администратор auth-proxy');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('nsinetskiy', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'nsinetskiy@rg.ru' , 'Никита Синецкий'  , 'Администратор auth-proxy');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('parshin', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'parshin@rg.ru' , 'Дмитрий Паршин'  , 'Администратор auth-proxy');
-- rosgas2011 => '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f'
-- 1 => '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b'


INSERT INTO app_user_role (appname, username, rolename) VALUES ('auth', 'admin', 'authadmin');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'admin', 'manager');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'admin', 'worker');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node1', 'admin', 'boss');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('node2', 'admin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('node2', 'admin', 'user');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'admin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'test', 'editor');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'chagin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'chernyshev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'ismagilov', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'kondratiev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'ivlev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'barsuk', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'filatchev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'nsinetskiy', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgcore', 'parshin', 'admin');


-- EXCEPTION WHEN OTHERS THEN 
--     RAISE EXCEPTION 'Тестовые данные уже существуют.';
-- END;
-- $$;
-- COMMIT;