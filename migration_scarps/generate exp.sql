SELECT DISTINCT a.appname, u.username, (a.appname || '_editor' ) AS rolename, random() AS r,
('INSERT INTO app_user_role (appname, username, rolename) VALUES (''' || 
a.appname || ''', ''' || u.username || ''', ''' || (a.appname || '_editor' ) || ''');' ) AS expr
FROM app AS a CROSS JOIN "user" AS u 
ORDER BY r
LIMIT 1500;
