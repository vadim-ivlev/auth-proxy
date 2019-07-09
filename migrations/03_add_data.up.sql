
DO 
$$
BEGIN
-- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

SET CONSTRAINTS app_user_role_fk_u DEFERRED;
SET CONSTRAINTS app_user_role_fk_a DEFERRED;

-- данные

INSERT INTO app (appname, description) VALUES ('app1'          , 'Тестовое express приложение 1 test0');
INSERT INTO app (appname, description) VALUES ('app2'          , 'Тестовое express приложение 2 test0');
INSERT INTO app (appname, description) VALUES ('onlinebc_admin', 'Тестовые трансляции test0'          );
INSERT INTO app (appname, description) VALUES ('rg'            , 'Сайт rg.ru test0'                   );

INSERT INTO app (appname, description) VALUES ('app10', 'phpMemcachedAdmin test');
INSERT INTO app (appname, description) VALUES ('app11', 'Push-уведомления test');
INSERT INTO app (appname, description) VALUES ('app12', 'Авторизация test');
INSERT INTO app (appname, description) VALUES ('app13', 'Акценты и блоки test');
INSERT INTO app (appname, description) VALUES ('app14', 'Баннерка test');
INSERT INTO app (appname, description) VALUES ('app15', 'Выставленное на сайт test');
INSERT INTO app (appname, description) VALUES ('app16', 'Документация test');
INSERT INTO app (appname, description) VALUES ('app17', 'Дубль два test');
INSERT INTO app (appname, description) VALUES ('app18', 'Задать вопрос test');
INSERT INTO app (appname, description) VALUES ('app19', 'Коды для вставки test');
INSERT INTO app (appname, description) VALUES ('app20', 'Комментарии test');
INSERT INTO app (appname, description) VALUES ('app21', 'Легкая версия сайта test');
INSERT INTO app (appname, description) VALUES ('app22', 'Магазин бумажной подписки (NEW) test');
INSERT INTO app (appname, description) VALUES ('app23', 'Медали test');
INSERT INTO app (appname, description) VALUES ('app24', 'Медиацентр test');
INSERT INTO app (appname, description) VALUES ('app25', 'Олимпиада 2014 test');
INSERT INTO app (appname, description) VALUES ('app26', 'Онлайн-трансляции test');
INSERT INTO app (appname, description) VALUES ('app27', 'Парсер test');
INSERT INTO app (appname, description) VALUES ('app28', 'Постинг в twitter test');
INSERT INTO app (appname, description) VALUES ('app29', 'Привязка регионов test');
INSERT INTO app (appname, description) VALUES ('app30', 'Родина test');
INSERT INTO app (appname, description) VALUES ('app31', 'Спортивные онлайн-трансляции test');
INSERT INTO app (appname, description) VALUES ('app32', 'Спортсмены test');
INSERT INTO app (appname, description) VALUES ('app33', 'Стена test');
INSERT INTO app (appname, description) VALUES ('app34', 'Теги test');
INSERT INTO app (appname, description) VALUES ('app35', 'Требуется на сайт test');
INSERT INTO app (appname, description) VALUES ('app36', 'Управление рассылками test');
INSERT INTO app (appname, description) VALUES ('app37', 'Фоторепортажи test');
INSERT INTO app (appname, description) VALUES ('app38', 'Чемпионаты test');




INSERT INTO "user" (username, password, email, fullname, description) VALUES ('vadim', '1', 'ivlev@rg.ru' , 'Ивлев Вадим'  , 'разработчик');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('serg' , '1', 'barsuk@rg.ru', 'Барсук Сергей', 'разработчик');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('max'  , '1', 'chagin@rg.ru', 'Чагин Максим' , 'начальник отдела разработки');

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
