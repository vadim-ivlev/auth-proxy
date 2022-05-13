
-- DO $$;
-- BEGIN
-- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

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
(6,'video'                 ,'http://host.docker.internal:7700'   , 'Видео GraphQL API админки'                                         ,NULL   ,NULL),
(7,'admin-comment'         ,'http://host.docker.internal:8095'   , 'GraphQL API комментарией'                                          ,NULL   ,NULL),
(8,'admin-core'            ,'http://host.docker.internal:8076'   , 'GraphQL API редакторского интерфейса'                              ,NULL   ,NULL),
(9,'file-uploader'         ,'http://host.docker.internal:8077'   , 'GraphQL API загрузки файлов'                                       ,NULL   ,NULL),
(10,'import'               ,'http://host.docker.internal:9099'   , 'GraphQL API импорт материалов в редактуру'                         ,NULL   ,NULL);

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные app уже существуют.';
END $$;


DO $$
BEGIN

INSERT INTO "user" 
(id, username    , password                                                          , email              , fullname           , description)
VALUES 
(1, 'admin'     , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'admin@rg.ru'      , 'Админ Админов'    , 'Администратор auth-proxy'),
(2, 'regdfo'    , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'dev1@rg.ru'       , 'Иван Иванов'     , 'Региональный редактор ДФО'         ),
(3, 'regcfo'    , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'dev2@rg.ru'       , 'Иван Иванов'     , 'Региональный редактор ЦФО'         ),
(4, 'chagin'    , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'chagin@rg.ru'     , 'Максим Чагин'     , 'Администратор auth-proxy'),
(5, 'chernyshev', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'chernyshev@rg.ru' , 'Алексей Чернышев' , 'Администратор auth-proxy'),
(6, 'boev'      , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'boev@rg.ru'       , 'Александр Боев'   , 'Администратор auth-proxy'),
(7, 'kondratiev', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'kondratiev@rg.ru' , 'Антон Кондратьев' , 'Администратор auth-proxy'),
(8, 'ivlev'     , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'ivlev@rg.ru'      , 'Вадим Ивлев'      , 'Администратор auth-proxy'),
(9, 'nsinetskiy', '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'nsinetskiy@rg.ru' , 'Никита Синецкий'  , 'Администратор auth-proxy'),
(10, 'pologov'   , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'pologov@rg.ru'    , 'Глеб Пологов'     , 'Администратор auth-proxy'),
(11, 'kataev'    , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'kataev@rg.ru'     , 'Антон Катаев'     , 'Администратор auth-proxy');

-- rosgas2011 => '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f'
-- 1 => '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b'

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные user уже существуют.';
END $$;



DO $$
BEGIN

INSERT INTO app_user_role 
(appname                 , username    , rolename)
VALUES 
('auth'                  , 'admin'     , 'authadmin'),

('echo'                  , 'admin'     , 'manager'),
('echo'                  , 'admin'     , 'worker'),
('echo'                  , 'admin'     , 'boss');

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные app_user_role уже существуют.';
END $$;


-- Группы ---------------------------------------------------------
DO $$
BEGIN

INSERT INTO "group"
(id, groupname       , description)
VALUES 
(1 ,'admins'         , 'Группа администраторов'),
(2 ,'guests'         , 'Посетители')            ,
(3 ,'developers'     , 'Разработчики'),
(4 ,'editors'        , 'Редакторы'),
(5 ,'comeditors'     , 'Выпускающие редакторы'),
(6 ,'regeditorsdfo'  , 'Региональные редакторы ДФО'),
(7 ,'regeditorscfo'  , 'Региональные редакторы ЦФО'),
(8 ,'bildeditors'    , 'Разработчики');

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные group уже существуют.';
END $$;




DO $$
BEGIN

INSERT INTO group_app_role 
(group_id, app_id, rolename)
VALUES 
(1, 6, 'admin'),
(1, 8, 'admin'),
(1, 10, 'admin'),
(4, 6, 'editor'),
(4, 8, 'editor'),
(4, 10, 'editor'),
(6, 6, 'viewer'),
(6, 8, 'regeditor'),
(6, 8, 'dfo'),
(7, 6, 'viewer'),
(7, 8, 'regeditor'),
(7, 8, 'cfo');

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные group_app_role уже существуют.';
END $$;


DO $$
BEGIN

INSERT INTO group_user_role 
(group_id, user_id)
VALUES 
(1, 1),
(1, 4),
(1, 5),
(1, 6),
(1, 7),
(1, 8),
(1, 9),
(1, 10),
(1, 11),
(6, 2),
(7, 3);

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные group_user_role уже существуют.';
END $$;


-- EXCEPTION WHEN OTHERS THEN 
--     RAISE EXCEPTION 'Тестовые данные уже существуют.';
-- END$$;
-- COMMIT;
