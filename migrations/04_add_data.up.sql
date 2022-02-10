
-- DO 
-- $$
-- BEGIN;
-- -- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

-- SET CONSTRAINTS app_user_role_fk_u DEFERRED;
-- SET CONSTRAINTS app_user_role_fk_a DEFERRED;

-- данные

DO $$
BEGIN

INSERT INTO app 
(id, appname                 , url                                 , description                                                         , public, sign)
VALUES
(0,'auth'                  ,''                                   , 'Сервис авторизации'                                                ,NULL   ,NULL),
(1,'echo'                  ,'https://echo-request.vercel.app/api', 'Показывает заголовки и тело запросов пропущенных через auth-proxy.',NULL   ,'Y' ),
(2,'echo-public'           ,'https://echo-request.vercel.app/api', 'Показывает заголовки и тело запросов пропущенных через auth-proxy.','Y'    ,NULL),
(3,'rg'                    ,'https://rg.ru'                      , 'Прокси к https://rg.ru'                                            ,'Y'    ,NULL),
(4,'photoreports-admin'    ,'http://host.docker.internal:8091'   , 'GraphQL API админки фоторепов'                                     ,NULL   ,NULL),
(5,'photoreports-admin-new','http://host.docker.internal:8094'   , 'Новое GraphQL API админки фоторепов'                               ,NULL   ,NULL),
(6,'video'                 ,'http://host.docker.internal:7700'   , 'Видео GraphQL API админки'                                         ,NULL   ,NULL),
(7,'admin-comment'         ,'http://host.docker.internal:8095'   , 'GraphQL API комментарией'                                          ,NULL   ,NULL),
(8,'rgcore'                ,'http://host.docker.internal:8076'   , 'GraphQL API редакторского интерфейса'                              ,NULL   ,NULL),
(9,'rgru-file-uploader'    ,'http://host.docker.internal:8077'   , 'GraphQL API загрузки файлов'                                       ,NULL   ,NULL);

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные app уже существуют.';
END $$;



DO $$
BEGIN

INSERT INTO "user" 
(username    , password                                                          , email              , fullname           , description)
VALUES 
('admin'     , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'admin@rg.ru'      , 'Админ Админов'    , 'Администратор auth-proxy'),
('test'      , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'dev@rg.ru'        , 'Иван Иванов'      , 'Tест auth-proxy'         ),
('chagin'    , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'chagin@rg.ru'     , 'Максим Чагин'     , 'Администратор auth-proxy'),
('chernyshev', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'chernyshev@rg.ru' , 'Алексей Чернышев' , 'Администратор auth-proxy'),
('boev'      , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'boev@rg.ru'       , 'Александр Боев'   , 'Администратор auth-proxy'),
('kondratiev', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'kondratiev@rg.ru' , 'Антон Кондратьев' , 'Администратор auth-proxy'),
('ivlev'     , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'ivlev@rg.ru'      , 'Вадим Ивлев'      , 'Администратор auth-proxy'),
('barsuk'    , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'barsuk@rg.ru'     , 'Сергей Барсук'    , 'Администратор auth-proxy'),
('filatchev' , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'filatchev@rg.ru'  , 'Дмитрий Филатчев' , 'Администратор auth-proxy'),
('nsinetskiy', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'nsinetskiy@rg.ru' , 'Никита Синецкий'  , 'Администратор auth-proxy'),
('pologov'   , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'pologov@rg.ru'    , 'Глеб Пологов'     , 'Администратор auth-proxy'),
('kataev'    , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'kataev@rg.ru'     , 'Антон Катаев'     , 'Администратор auth-proxy');

-- rosgas2011 => '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f'
-- 1 => '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b'

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные user уже существуют.';
END $$;


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
('photoreports-admin-new', 'boev'      , 'admin'),
('photoreports-admin-new', 'kondratiev', 'admin'),
('photoreports-admin-new', 'ivlev'     , 'admin'),
('photoreports-admin-new', 'barsuk'    , 'admin'),
('photoreports-admin-new', 'filatchev' , 'admin'),
('photoreports-admin-new', 'nsinetskiy', 'admin'),
('photoreports-admin-new', 'kataev'    , 'admin'),
('photoreports-admin-new', 'pologov'   , 'admin'),

('rgcore'                , 'admin'     , 'admin'),
('rgcore'                , 'test'      , 'editor'),
('rgcore'                , 'chagin'    , 'admin'),
('rgcore'                , 'chernyshev', 'admin'),
('rgcore'                , 'boev'      , 'admin'),
('rgcore'                , 'kondratiev', 'admin'),
('rgcore'                , 'ivlev'     , 'admin'),
('rgcore'                , 'barsuk'    , 'admin'),
('rgcore'                , 'filatchev' , 'admin'),
('rgcore'                , 'nsinetskiy', 'admin'),
('rgcore'                , 'kataev'    , 'admin'),
('rgcore'                , 'pologov'   , 'admin'),

('rgru-file-uploader'    , 'admin'     , 'admin'),
('rgru-file-uploader'    , 'test'      , 'editor'),
('rgru-file-uploader'    , 'chagin'    , 'admin'),
('rgru-file-uploader'    , 'chernyshev', 'admin'),
('rgru-file-uploader'    , 'boev'      , 'admin'),
('rgru-file-uploader'    , 'kondratiev', 'admin'),
('rgru-file-uploader'    , 'ivlev'     , 'admin'),
('rgru-file-uploader'    , 'barsuk'    , 'admin'),
('rgru-file-uploader'    , 'filatchev' , 'admin'),
('rgru-file-uploader'    , 'nsinetskiy', 'admin'),
('rgru-file-uploader'    , 'kataev'    , 'admin'),
('rgru-file-uploader'    , 'pologov'   , 'admin'),

('video'                 , 'admin'     , 'admin'),
('video'                 , 'test'      , 'editor'),
('video'                 , 'chagin'    , 'admin'),
('video'                 , 'chernyshev', 'admin'),
('video'                 , 'boev'      , 'admin'),
('video'                 , 'kondratiev', 'admin'),
('video'                 , 'ivlev'     , 'admin'),
('video'                 , 'barsuk'    , 'admin'),
('video'                 , 'filatchev' , 'admin'),
('video'                 , 'nsinetskiy', 'admin'),
('video'                 , 'kataev'    , 'admin'),
('video'                 , 'pologov'   , 'admin'),

('admin-comment'         , 'admin'     , 'admin'),
('admin-comment'         , 'test'      , 'editor'),
('admin-comment'         , 'chagin'    , 'admin'),
('admin-comment'         , 'chernyshev', 'admin'),
('admin-comment'         , 'boev'      , 'admin'),
('admin-comment'         , 'kondratiev', 'admin'),
('admin-comment'         , 'ivlev'     , 'admin'),
('admin-comment'         , 'barsuk'    , 'admin'),
('admin-comment'         , 'filatchev' , 'admin'),
('admin-comment'         , 'nsinetskiy', 'admin'),
('admin-comment'         , 'kataev'    , 'admin'),
('admin-comment'         , 'pologov'   , 'admin');



-- Группы ---------------------------------------------------------

INSERT INTO "group"
(id, groupname       , description)
VALUES 
(1 ,'administrators' , 'Группа администраторов'),
(2 ,'guests'         , 'Посетители')            ,
(3 ,'developers'     , 'Разработчики');

INSERT INTO group_app_role 
(group_id, app_id, rolename)
VALUES 
(1, 0, 'authadmin'),
(2, 1, 'guestrole1'),
(2, 1, 'guestrole2');


-- EXCEPTION WHEN OTHERS THEN 
--     RAISE EXCEPTION 'Тестовые данные уже существуют.';
-- END;
-- $$;
-- COMMIT;
