
-- DO 
-- $$
-- BEGIN;
-- -- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

-- SET CONSTRAINTS app_user_role_fk_u DEFERRED;
-- SET CONSTRAINTS app_user_role_fk_a DEFERRED;

-- данные


INSERT INTO app (appname, url, description) VALUES ('auth','', 'Сервис авторизации');
INSERT INTO app (appname, url, description, sign) VALUES ('echo','https://echo-request.vercel.app/api', 'Показывает заголовки и тело запросов пропущенных через auth-proxy.', 'Y');
INSERT INTO app (appname, url, description, public) VALUES ('rg','https://rg.ru', 'Прокси к https://rg.ru', 'Y' );
INSERT INTO app (appname, url, description) VALUES ('photoreports-admin','http://host.docker.internal:8091', 'GraphQL API админки фоторепов');
INSERT INTO app (appname, url, description) VALUES ('photoreports-admin-new','http://host.docker.internal:8094', 'Новое GraphQL API админки фоторепов');
INSERT INTO app (appname, url, description) VALUES ('video','http://host.docker.internal:7700', 'Видео GraphQL API админки');
INSERT INTO app (appname, url, description) VALUES ('rgcore','http://host.docker.internal:8076', 'GraphQL API редакторского интерфейса');
INSERT INTO app (appname, url, description) VALUES ('rgru-file-uploader','http://host.docker.internal:8077', 'GraphQL API загрузки файлов');


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

INSERT INTO app_user_role (appname, username, rolename) VALUES ('echo', 'admin', 'manager');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('echo', 'admin', 'worker');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('echo', 'admin', 'boss');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin', 'admin', 'admin');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'admin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'test', 'editor');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'chagin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'chernyshev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'ismagilov', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'kondratiev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'ivlev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'barsuk', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'filatchev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'nsinetskiy', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('photoreports-admin-new', 'parshin', 'admin');

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

INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'admin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'test', 'editor');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'chagin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'chernyshev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'ismagilov', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'kondratiev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'ivlev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'barsuk', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'filatchev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'nsinetskiy', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('rgru-file-uploader', 'parshin', 'admin');

INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'admin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'test', 'editor');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'chagin', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'chernyshev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'ismagilov', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'kondratiev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'ivlev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'barsuk', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'filatchev', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'nsinetskiy', 'admin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('video', 'parshin', 'admin');


-- EXCEPTION WHEN OTHERS THEN 
--     RAISE EXCEPTION 'Тестовые данные уже существуют.';
-- END;
-- $$;
-- COMMIT;
