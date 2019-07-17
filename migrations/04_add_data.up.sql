
-- DO 
-- $$
BEGIN;
-- -- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

-- SET CONSTRAINTS app_user_role_fk_u DEFERRED;
-- SET CONSTRAINTS app_user_role_fk_a DEFERRED;

-- данные

INSERT INTO app (appname, url, description) VALUES ('auth'          ,'', 'Сервис авторизации');
INSERT INTO app (appname, url, description) VALUES ('app1'          ,'http://node:3001', 'Express приложение 1. Работает в докере. test0');
INSERT INTO app (appname, url, description) VALUES ('app2'          ,'http://localhost:3002', 'Тестовое express приложение 2. test0');
INSERT INTO app (appname, url, description) VALUES ('onlinebc_admin','http://localhost:7700', 'Тестовые трансляции test0'          );
INSERT INTO app (appname, url, description) VALUES ('rg'            ,'https://rg.ru'        , 'Сайт rg.ru test0'                   );

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




INSERT INTO "user" (username, password, email, fullname, description) VALUES ('vadim', '1', 'ivlev@rg.ru' , 'Ивлев Вадим'  , 'разработчик test0');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('serg' , '1', 'barsuk@rg.ru', 'Барсук Сергей', 'разработчик test0');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('max'  , '1', 'chagin@rg.ru', 'Чагин Максим' , 'начальник отдела разработки test0');

INSERT INTO "user" (username, password, email, fullname, description) VALUES ('7101159', 'q', '7101159@rambler.ru', 'Журман Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('a.bondarev', 'q', 'a.bondarev@krasrg.ru', 'Бондарев Алексей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('a.sorokina', 'q', 'a.sorokina@rbth.ru', 'Sorokina Anna ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('achernega4', 'q', 'achernega4@rg.ru', 'Чернега Александра', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('afinez', 'q', 'afinez@mail.ru', 'Карнаухов Игорь', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('agur', 'q', 'agur@rg.ru', 'Жанна Агуреева', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('akondratiev', 'q', 'akondratiev@rg.ru', 'Кондратьев Антон', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('akulova', 'q', 'akulova@rg.ru', 'Акулова Евгения', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('aleksanyan', 'q', 'aleksanyan@rg.ru', 'Алексанян Элен', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('alex', 'q', 'alex@rg.nsk.su', 'Хадаев Алексей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('alex3755', 'q', 'alex3755@mail.ru', 'Воздвиженская Александра ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('aminova', 'q', 'aminova@rg.ru', 'Аминова Марина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('andreyandreev1961', 'q', 'andreyandreev1961@yandex.ru', 'Андреев Андрей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('anna_alex_', 'q', 'anna_alex_@mail.ru', 'Бондаренко Анна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('anna-skripka', 'q', 'anna-skripka@mail.ru', 'Скрипка Анна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('annafedor2009', 'q', 'annafedor2009@yandex.ru', 'Шепелева Анна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('annaskudayeva', 'q', 'annaskudayeva@gmail.com', 'Скудаева Анна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('antivalagin', 'q', 'antivalagin@yandex.ru', 'Валагин Антон', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('any888', 'q', 'any888@mail.ru', 'Тимофеева Анна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('apol3lp', 'q', 'apol3lp@mail.ru', 'Полухин Андрей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ashirova', 'q', 'ashirova@rg.ru', 'Аширова Эльмира', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('atan', 'q', 'atan@rg-ural.ru', 'Андреева Татьяна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('barsuk', 'q', 'barsuk@rg.ru', 'Барсук Сергей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('Barybina', 'q', 'Barybina@rg.ru', 'Барыбина Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('bbkilin', 'q', 'bbkilin@gmail.com', 'Килин Борис', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('bil11', 'q', 'bil11@mail.ru', 'Биль Владимир', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('bormasheva', 'q', 'bormasheva@rg.ru', 'Бормашева Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('borodenko', 'q', 'borodenko@rg.ru', 'Бороденко Виктория', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('brook', 'q', 'brook@rg.ru', 'Брук Елена', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('burlo', 'q', 'burlo@rg.ru', 'Бурло Алексей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('busuek', 'q', 'busuek@rg.ru', 'Бусуек Светлана ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('butylkina', 'q', 'butylkina@rg.ru', 'Бутылкина Нина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('buynaya', 'q', 'buynaya@rg.ru', 'Буйная Юлия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('carminesaviano', 'q', 'carminesaviano@gmail.com', 'Zavialov Vladimir ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('chagin', 'q', 'chagin@rg.ru', 'Чагин Максим', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('chernyakov', 'q', 'chernyakov@rg.ru', 'Черняков Евгений', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('chernyshev', 'q', 'chernyshev@rg.ru', 'Чернышев Алексей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('chipak', 'q', 'chipak@rg.ru', 'Чипак Павел ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('chulkov', 'q', 'chulkov@rg.ru', 'Чулков Дмитрий ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('david.flores', 'q', 'david.flores@rbth.ru', 'Flores David ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('deeva', 'q', 'deeva@rg.ru', 'Деева Наталья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('den77p', 'q', 'den77p@mail.ru', 'Передельский Денис', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('denissamara', 'q', 'denissamara@yandex.ru', 'Кудряшов Денис ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('dgulnaz', 'q', 'dgulnaz@yandex.ru', 'Данилова Гульназ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('dimagrigor58', 'q', 'dimagrigor58@gmail.com', 'Григорьев Дмитрий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('dinessa', 'q', 'dinessa@hotmail.com', 'Доценко Инесса ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('dtverdyy', 'q', 'dtverdyy@rg.ru', 'Твердый Дмитрий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('dunay.es', 'q', 'dunay.es@gmail.com', 'Хайновская Евгения', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('e.savchenko743', 'q', 'e.savchenko743@gmail.com', 'Савченко Екатерина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('economy', 'q', 'economy@rg.nnov.ru', 'Норский Валерий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('elias13', 'q', 'elias13@yandex.ru', 'Соболев Илья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('elisaveta_golub', 'q', 'elisaveta_golub@mail.ru', 'Голуб Лиза', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ershov', 'q', 'ershov@rgkuban.ru', 'Ершов Александр', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('esolovyev', 'q', 'esolovyev@rg.ru', 'Соловьев Ефим ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('evolkovy', 'q', 'evolkovy@gmail.com', 'Волков Евгений ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('filimonov', 'q', 'filimonov@rg.ru', 'Филимонов Дмитрий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('filippov2', 'q', 'filippov2@rg.ru', 'Филиппов Илья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('filonofz', 'q', 'filonofz@mail.ru', 'Филонов Игорь', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('FKazakovK', 'q', 'FKazakovK@yandex.ru', 'Казаков Федор', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('fofana', 'q', 'fofana@mail.ru', 'Ткачева Татьяна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('front09', 'q', 'front09@rambler.ru', 'Федосов Александр', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('gavrilenko-apk', 'q', 'gavrilenko-apk@mail.ru', 'Гавриленко Александр', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('gerasimenkoe', 'q', 'gerasimenkoe@rg.ru', 'Герасименко Екатерина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('glushenkov', 'q', 'glushenkov@rg.ru', 'Glushenkov Yuri ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('gomozov', 'q', 'gomozov@rg.ru', 'Гомозов Виктор ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('grafn', 'q', 'grafn@inbox.ru', 'Граф Наталья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('grishenko', 'q', 'grishenko@rostov.rg.ru', 'Грищенко Николай', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('hajoff', 'q', 'hajoff@mail.ru', 'Орлов Сергей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('idrobysheva', 'q', 'idrobysheva@yandex.ru', 'Дробышева Ирина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ifedorov', 'q', 'ifedorov@rg.ru', 'Федоров Игорь', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('il2009', 'q', 'il2009@list.ru', 'Изотов Илья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('inin', 'q', 'inin@rg.ru', 'Инин Дмитрий ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('iponosov', 'q', 'iponosov@yandex.ru', 'Поносов Илья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('irina_evsukova00', 'q', 'irina_evsukova00@mail.ru', 'Евсюкова Ирина ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ismagilov', 'q', 'ismagilov@rg.ru', 'Исмагилов Камиль', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ivanitskaya', 'q', 'ivanitskaya@rg.ru', 'Иваницкая Полина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ivlev', 'q', 'ivlev@rg.ru', 'Ивлев Вадим', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('izubko', 'q', 'izubko@list.ru', 'Зубко Илья ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('johnsylver', 'q', 'johnsylver@yandex.ru', 'Васильев Алексей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('julia_gardner_rg', 'q', 'julia_gardner_rg@mail.ru', 'Гарднер Юлия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kandy68', 'q', 'kandy68@list.ru', 'Куликов Андрей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kasya20', 'q', 'kasya20@yandex.ru', 'Дубичева Ксения', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kazans-oleg', 'q', 'kazans-oleg@yandex.ru', 'Корякин Олег', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kazmina', 'q', 'kazmina@rg.ru', 'Казьмина Елена', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kburmenko', 'q', 'kburmenko@yandex.ru', 'Бурменко Ксения ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kerdo', 'q', 'kerdo@yandex.ru', 'Бабкин Сергей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('khoroshilov', 'q', 'khoroshilov@rg.ru', 'Хорошилов Михаил', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kmerinov', 'q', 'kmerinov@rg.ru', 'Меринов Константин ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kolesnikov', 'q', 'kolesnikov@rg.ru', 'Колесников Вадим ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kolmykova', 'q', 'kolmykova@rg.ru', 'Колмыкова Наталья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('komarovaa', 'q', 'komarovaa@rg.ru', 'Комарова Анастасия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kondreva', 'q', 'kondreva@rgkazan.mi.ru', 'Кондрева Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kor', 'q', 'kor@rgkazan.mi.ru', 'Брайловская Светлана ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('koshkav', 'q', 'koshkav@yandex.ru', 'Чернышева Виктория', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kras', 'q', 'kras@rg.ru', 'Краснянская Виктория', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('krimjuly', 'q', 'krimjuly@gmail.com', 'Крымова Юлия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('kushnareva', 'q', 'kushnareva@rg.ru', 'Кушнарева Ася', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('lashkul', 'q', 'lashkul@rg.ru', 'Лашкул Никита', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('lenacarn2005', 'q', 'lenacarn2005@mail.ru', 'Исакова Елена', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('litovchenko', 'q', 'litovchenko@rg.ru', 'Литовченко Алексей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('lokalov', 'q', 'lokalov@rg.ru', 'Локалов Артем', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('lolasmit', 'q', 'lolasmit@mail.ru', 'Сутула Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('lucia.bellinello', 'q', 'lucia.bellinello@rbth.ru', 'Bellinello Lucia ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('m.pinkus', 'q', 'm.pinkus@mail.ru', 'Пинкус Михаил', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('malikova', 'q', 'malikova@rg.ru', 'Маликова Евгения', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('malinina', 'q', 'malinina@rg.ru', 'Малинина Анна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('malysheva', 'q', 'malysheva@rg.ru', 'Малышева Наталья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('man72t', 'q', 'man72t@mail.ru', 'Меньшиков Анатолий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('marshak', 'q', 'marshak@rg.ru', 'Маршак Илья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('maslova', 'q', 'maslova@rg.ru', 'Маслова Татьяна ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('matguka', 'q', 'matguka@mail.ru', 'Камаева Айгуль', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('maximov', 'q', 'maximov@rg.ru', 'Максимов Илья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('mediakat', 'q', 'mediakat@yandex.ru', 'Дементьева Екатерина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('melkumyan', 'q', 'melkumyan@rg.ru', 'Мелкумян Карина ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('melnik', 'q', 'melnik@rg.ru', 'Мельник Оксана', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('melnikov', 'q', 'melnikov@rg.ru', 'Мельников Андрей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('mil', 'q', 'mil@rg-ural.ru', 'Миляева Елена', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('mnovikov', 'q', 'mnovikov@rg.ru', 'Новиков Максим ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('mogana.ks', 'q', 'mogana.ks@gmail.com', 'Семенко Ксения ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('morozova_o', 'q', 'morozova_o@rg.ru', 'Ольга Морозова', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('morrisspb', 'q', 'morrisspb@gmail.com', 'Голубкова Мария', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('nad8417', 'q', 'nad8417@yandex.ru', 'Столярчук Надежда', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('natali-shika', 'q', 'natali-shika@yandex.ru', 'Широкова Наталья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('nikiforova', 'q', 'nikiforova@rg.ru', 'Никифорова Полина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('nikolavna2-2', 'q', 'nikolavna2-2@yandex.ru', 'Грибанова Оксана', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('novik.artem', 'q', 'novik.artem@gmail.com', 'Новиков Артем ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('nsinetskiy', 'q', 'nsinetskiy@rg.ru', 'Синецкий Никита', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('nurimanova', 'q', 'nurimanova@rg.ru', 'Нуриманова Зульфия ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('obiedkova', 'q', 'obiedkova@rg.ru', 'Объедкова Елена', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('olle73', 'q', 'olle73@mail.ru', 'Горюнова Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('osilbekova.d', 'q', 'osilbekova@inbox.ru', 'Осильбекова Дина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('osilbekova.m', 'q', 'osilbekova@rg.ru', 'Осильбекова Мария', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ovarlamova', 'q', 'ovarlamova@yandex.ru', 'Дмитренко Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('pan', 'q', 'pan@rg-ural.ru', 'Панасенко Сергей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('panin', 'q', 'panin@rg.ru', 'Панин Георгий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('parshin', 'q', 'parshin@rg.ru', 'Паршин Дмитрий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('pavlovam', 'q', 'pavlovam@rg.ru', 'Павлова Мария', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('peteshova', 'q', 'peteshova@rg.ru', 'Петешова Елена ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('petrovan', 'q', 'petrovan@rg.ru', 'Петров Алексей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('petyaeva', 'q', 'petyaeva@rg.ru', 'Петяева Анна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('pichurina45', 'q', 'pichurina45@mail.ru', 'Пичурина Валентина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('pkolga76', 'q', 'pkolga76@gmail.com', 'Бухарова Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('polina.kortina', 'q', 'polina.kortina@rbth.ru', 'Kortina Polina ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('pomeshchikova', 'q', 'pomeshchikova@rg.ru', 'Помещикова Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('pr2', 'q', 'pr2@rg-ural.ru', 'Швабауэр Наталия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('prasolov', 'q', 'prasolov@rg.ru', 'Прасолов Олег', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('promyshlyaeva', 'q', 'promyshlyaeva@rg.ru', 'Промышляева Олеся', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ptax66', 'q', 'ptax66@bk.ru', 'Войтович Ирина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ptr2002', 'q', 'ptr2002@mail.ru', 'Патрикеева Юлия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('pulya', 'q', 'pulya@rbth.ru', 'Pulya Vsevolod', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('pusto-pusto', 'q', 'pusto-pusto@list.ru', 'Субботин Георгий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('raichev', 'q', 'raichev@list.ru', 'Раичев Дмитрий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ramilya_tr', 'q', 'ramilya_tr@mail.ru', 'Туктарова Рамиля', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('razuvaeva', 'q', 'razuvaeva@rg.ru', 'Разуваева Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('redaktor2', 'q', 'redaktor2@rgkazan.mi.ru', 'Аристова Галина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('rg-saransk', 'q', 'rg-saransk@bk.ru', 'Зотикова Валентина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('rg-tver', 'q', 'rg-tver@yandex.ru', 'Кузнецов Денис', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('rgkaterina', 'q', 'rgkaterina@inbox.ru', 'Ковалевская Екатерина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('rgkazan-inet', 'q', 'rgkazan-inet@yandex.ru', 'Алимов Тимур', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('rgskfo', 'q', 'rgskfo@mail.ru', 'Брежицкая Елена ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('rgstav', 'q', 'rgstav@yandex.ru', 'Емельянова Светлана ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('rgural', 'q', 'rgural@rg-ural.ru', 'Резникова Юлия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('rkiselev', 'q', 'rkiselev@rbth.ru', 'Kiselev Roman ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('romrom1983', 'q', 'romrom1983@bk.ru', 'Романенко Роман', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('Rudavina', 'q', 'Rudavina@rg.ru', 'Рудавина Татьяна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('ruleva', 'q', 'ruleva@rg-ural.ru', 'Рулева Наталья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('s.sobolev', 'q', 's.sobolev@rg.ru', 'Соболев Сергей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('sadriev', 'q', 'sadriev@rg.ru', 'Садриев Ильмас', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('samar', 'q', 'samar@rg.ru', 'Самарская Мария', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('san', 'q', 'san@rg-ural.ru', 'Санатина Юлия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('satoyusuke', 'q', 'satoyusuke@list.ru', 'Yusuke Sato ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('sedkov', 'q', 'sedkov@rg.ru', 'Седьков Дмитрий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('shamilkerashev', 'q', 'shamilkerashev@gmail.com', 'Керашев Шамиль', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('shashaeva', 'q', 'shashaeva@rg.ru', 'Шашаева Мария', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('shestakov', 'q', 'shestakov@rg.perm.ru', 'Шестаков Александр', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('shu_lena', 'q', 'shu_lena@mail.ru', 'Шулепова Елена', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('simonov', 'q', 'simonov@rg.ru', 'Симонов Павел ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('site', 'q', 'site@rg-ural.ru', 'Воробьева Татьяна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('skorpy5', 'q', 'skorpy5@mail.ru', 'Мационг Елена', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('smeyukha', 'q', 'smeyukha@rg.ru', 'Смеюха Виктор ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('soldatenkov', 'q', 'soldatenkov@rg.ru', 'Солдатенков Максим', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('sosnovsky', 'q', 'sosnovsky@rg.ru', 'Сосновский Дмитрий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('suhorukova', 'q', 'suhorukova@rg.ru', 'Сухорукова Анастасия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('suvorova', 'q', 'suvorova@agima.ru', 'Суварова Юлия ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('svetarg', 'q', 'svetarg@yandex.ru', 'Песоцкая Светлана', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('svetl', 'q', 'svetl@rg-ural.ru', 'Добрынина Светлана', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('svetlana-tsygankova', 'q', 'svetlana-tsygankova@yandex.ru', 'Цыганкова Светлана', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('swan75', 'q', 'swan75@mail.ru', 'Дмитракова Татьяна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('tata010', 'q', 'tata010@yandex.ru', 'Саванкова Наталья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('tatiana-rg', 'q', 'tatiana-rg@yandex.ru', 'Ермошкина Татьяна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('tatyana.shramchenko', 'q', 'tatyana.shramchenko@rbth.ru', 'Shramchenko Tatiana ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('telyatnikova', 'q', 'telyatnikova@rg.ru', 'Телятникова Снежана', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('tenevaya', 'q', 'tenevaya@gmail.com', 'Чубарова Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('test', 'q', 'test@test.com', 'Тест Тест ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('tkolesnikova', 'q', 'tkolesnikova@rg.ru', 'Колесникова Татьяна ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('tovprapor', 'q', 'tovprapor@yandex.ru', 'Словохотов Сергей', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('trhn', 'q', 'trhn@mail.ru', 'Труханова Элина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('trir76', 'q', 'trir76@inbox.ru', 'Бикбаев Марсель', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('tromanova', 'q', 'tromanova@rg.ru', 'Романова Татьяна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('trubko', 'q', 'trubko@rg.ru', 'Трубко Мария', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('trusova', 'q', 'trusova@rg.ru', 'Трусова Анастасия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('tuktarova', 'q', 'tuktarova@rg.ru', 'Туктарова Рамиля', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('turkova.1', 'q', 'turkova.1@mail.ru', 'Туркова Татьяна ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('usov', 'q', 'usov@rg.ru', 'Усов Денис', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('uvylegzhanina', 'q', 'uvylegzhanina@yandex.ru', 'Вылегжанина Ульяна', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('v_ms', 'q', 'v_ms@mail.ru', 'Волкова Мария', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('v.cherenewa', 'q', 'v.cherenewa@yandex.ru', 'Черенева Вера', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('vadim_dav', 'q', 'vadim_dav@rambler.ru', 'Давыденко Вадим', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('valkichin', 'q', 'valkichin@mail.ru', 'Кичин Валерий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('varyukhina', 'q', 'varyukhina@rg.ru', 'Варюхина Екатерина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('vbpetrov', 'q', 'vbpetrov@mail.ru', 'Петров Владимир', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('victor', 'q', 'victor@rg.nnov.ru', 'Девицын Виктор', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('voronina-andreeva', 'q', 'voronina-andreeva@yandex.ru', 'Воронина Дарья', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('vovata', 'q', 'vovata@yandex.ru', 'Таюрский Владимир', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('vrynishka_ya', 'q', 'vrynishka_ya@mail.ru', 'Бондаренко Ольга', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('vs-anti', 'q', 'vs-anti@yandex.ru', 'Латухина Кира ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('vyunok79', 'q', 'vyunok79@mail.ru', 'Лобанова Екатерина', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('yuliyasyrova111', 'q', 'yuliyasyrova111@gmail.com', 'Колбина Юлия', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('yurchak', 'q', 'yurchak@rg.ru', 'Юрчак Виталий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('yzubko', 'q', 'yzubko@rg.ru', 'Зубко Юрий', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('zadunayskaya', 'q', 'zadunayskaya@rg.ru', 'Задунайская Оксана', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('zakiev', 'q', 'zakiev@rg.ru', 'Закиев Ренат', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('zakieva', 'q', 'zakieva@rg.ru', 'Закиева Гульнара ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('zhilyaev', 'q', 'zhilyaev@rg.ru', 'Жиляев Данила', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('zinchenko', 'q', 'zinchenko@rg.ru', 'Зинченко Алексей ', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('zinkler', 'q', 'zinkler@yandex.ru', 'Цинклер Евгения', 'test');
INSERT INTO "user" (username, password, email, fullname, description) VALUES ('zzzebra', 'q', 'zzzebra@gmail.com', 'Благовещенский Антон', 'test');


INSERT INTO app_user_role (appname, username, rolename) VALUES ('auth', 'vadim', 'authadmin');
INSERT INTO app_user_role (appname, username, rolename) VALUES ('auth', 'vadim', 'manager');

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


-- EXCEPTION WHEN OTHERS THEN 
--     RAISE EXCEPTION 'Тестовые данные уже существуют.';
-- END;
-- $$;
COMMIT;
