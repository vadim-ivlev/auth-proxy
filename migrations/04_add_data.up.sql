
-- DO 
-- $$
-- BEGIN;
-- -- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

-- SET CONSTRAINTS app_user_role_fk_u DEFERRED;
-- SET CONSTRAINTS app_user_role_fk_a DEFERRED;

-- данные



INSERT INTO app 
(appname                 , url                                 , description                                                         , public, sign)
VALUES
('auth'                  ,''                                   , 'Сервис авторизации'                                                ,NULL   ,NULL),
('echo'                  ,'https://echo-request.vercel.app/api', 'Показывает заголовки и тело запросов пропущенных через auth-proxy.',NULL   ,'Y' ),
('echo-public'           ,'https://echo-request.vercel.app/api', 'Показывает заголовки и тело запросов пропущенных через auth-proxy.','Y'    ,NULL),
('rg'                    ,'https://rg.ru'                      , 'Прокси к https://rg.ru'                                            ,'Y'    ,NULL),
('photoreports-admin'    ,'http://host.docker.internal:8091'   , 'GraphQL API админки фоторепов'                                     ,NULL   ,NULL),
('photoreports-admin-new','http://host.docker.internal:8094'   , 'Новое GraphQL API админки фоторепов'                               ,NULL   ,NULL),
('video'                 ,'http://host.docker.internal:7700'   , 'Видео GraphQL API админки'                                         ,NULL   ,NULL),
('rgcore'                ,'http://host.docker.internal:8076'   , 'GraphQL API редакторского интерфейса'                              ,NULL   ,NULL),
('rgru-file-uploader'    ,'http://host.docker.internal:8077'   , 'GraphQL API загрузки файлов'                                       ,NULL   ,NULL);




INSERT INTO "user" 
(username    , password                                                          , email              , fullname           , description)
VALUES 
('admin'     , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'admin@rg.ru'      , 'Админ Админов'    , 'Администратор auth-proxy'),
('test'      , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'dev@rg.ru'        , 'Иван Иванов'      , 'Tест auth-proxy'         ),
('chagin'    , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'chagin@rg.ru'     , 'Максим Чагин'     , 'Администратор auth-proxy'),
('chernyshev', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'chernyshev@rg.ru' , 'Алексей Чернышев' , 'Администратор auth-proxy'),
('ismagilov' , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'ismagilov@rg.ru'  , 'Камиль Исмагилов' , 'Администратор auth-proxy'),
('kondratiev', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'kondratiev@rg.ru' , 'Антон Кондратьев' , 'Администратор auth-proxy'),
('ivlev'     , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'ivlev@rg.ru'      , 'Вадим Ивлев'      , 'Администратор auth-proxy'),
('barsuk'    , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'barsuk@rg.ru'     , 'Сергей Барсук'    , 'Администратор auth-proxy'),
('filatchev' , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'filatchev@rg.ru'  , 'Дмитрий Филатчев' , 'Администратор auth-proxy'),
('nsinetskiy', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'nsinetskiy@rg.ru' , 'Никита Синецкий'  , 'Администратор auth-proxy'),
('parshin'   , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'parshin@rg.ru'    , 'Дмитрий Паршин'   , 'Администратор auth-proxy');

-- rosgas2011 => '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f'
-- 1 => '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b'


INSERT INTO app_user_role 
(appname                 , username    , rolename)
VALUES 
('auth'                  , 'admin'     , 'authadmin'),

('echo'                  , 'admin'     , 'manager'),
('echo'                  , 'admin'     , 'worker'),
('echo'                  , 'admin'     , 'boss'),

('photoreports-admin'    , 'admin'     , 'admin'),
('photoreports-admin-new', 'admin'     , 'admin'),
('photoreports-admin-new', 'test'      , 'editor'),
('photoreports-admin-new', 'chagin'    , 'admin'),
('photoreports-admin-new', 'chernyshev', 'admin'),
('photoreports-admin-new', 'ismagilov' , 'admin'),
('photoreports-admin-new', 'kondratiev', 'admin'),
('photoreports-admin-new', 'ivlev'     , 'admin'),
('photoreports-admin-new', 'barsuk'    , 'admin'),
('photoreports-admin-new', 'filatchev' , 'admin'),
('photoreports-admin-new', 'nsinetskiy', 'admin'),
('photoreports-admin-new', 'parshin'   , 'admin'),

('rgcore'                , 'admin'     , 'admin'),
('rgcore'                , 'test'      , 'editor'),
('rgcore'                , 'chagin'    , 'admin'),
('rgcore'                , 'chernyshev', 'admin'),
('rgcore'                , 'ismagilov' , 'admin'),
('rgcore'                , 'kondratiev', 'admin'),
('rgcore'                , 'ivlev'     , 'admin'),
('rgcore'                , 'barsuk'    , 'admin'),
('rgcore'                , 'filatchev' , 'admin'),
('rgcore'                , 'nsinetskiy', 'admin'),
('rgcore'                , 'parshin'   , 'admin'),

('rgru-file-uploader'    , 'admin'     , 'admin'),
('rgru-file-uploader'    , 'test'      , 'editor'),
('rgru-file-uploader'    , 'chagin'    , 'admin'),
('rgru-file-uploader'    , 'chernyshev', 'admin'),
('rgru-file-uploader'    , 'ismagilov' , 'admin'),
('rgru-file-uploader'    , 'kondratiev', 'admin'),
('rgru-file-uploader'    , 'ivlev'     , 'admin'),
('rgru-file-uploader'    , 'barsuk'    , 'admin'),
('rgru-file-uploader'    , 'filatchev' , 'admin'),
('rgru-file-uploader'    , 'nsinetskiy', 'admin'),
('rgru-file-uploader'    , 'parshin'   , 'admin'),

('video'                 , 'admin'     , 'admin'),
('video'                 , 'test'      , 'editor'),
('video'                 , 'chagin'    , 'admin'),
('video'                 , 'chernyshev', 'admin'),
('video'                 , 'ismagilov' , 'admin'),
('video'                 , 'kondratiev', 'admin'),
('video'                 , 'ivlev'     , 'admin'),
('video'                 , 'barsuk'    , 'admin'),
('video'                 , 'filatchev' , 'admin'),
('video'                 , 'nsinetskiy', 'admin'),
('video'                 , 'parshin'   , 'admin');


-- EXCEPTION WHEN OTHERS THEN 
--     RAISE EXCEPTION 'Тестовые данные уже существуют.';
-- END;
-- $$;
-- COMMIT;
