SELECT DISTINCT a.appname, u.username, (a.appname || '_editor' ) AS rolename, random() AS r,
('INSERT INTO app_user_role (appname, username, rolename) VALUES (''' || 
a.appname || ''', ''' || u.username || ''', ''' || (a.appname || '_editor' ) || ''');' ) AS expr
FROM app AS a CROSS JOIN "user" AS u 
ORDER BY r
LIMIT 1500;



INSERT INTO app (appname, url, description)         VALUES ('auth'      ,'', 'Сервис авторизации');
INSERT INTO app (appname, url, description)         VALUES ('app1'      ,'http://node:3001', 'Express приложение 1. Работает в докере. test0');
INSERT INTO app (appname, url, description)         VALUES ('app2'      ,'http://localhost:3002', 'Тестовое express приложение 2. test0');
INSERT INTO app (appname, url, description, rebase) VALUES ('pravo'     ,'https://pravo.rg.ru', 'Право test0' ,'Y' );
INSERT INTO app (appname, url, description)         VALUES ('rg'        ,'https://rg.ru'      , 'Сайт rg.ru test0' );
INSERT INTO app (appname, url, description)         VALUES ('wikipedia' ,'https://www.wikipedia.org', 'Википедия test0');
