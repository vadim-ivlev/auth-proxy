
-- DO $$;
-- BEGIN
-- чтобы не заботиться о порядке вставки данных из за ограничений внешних ключей

-- SET CONSTRAINTS app_user_role_fk_u DEFERRED;
-- SET CONSTRAINTS app_user_role_fk_a DEFERRED;

-- данные

DO $$
BEGIN

	INSERT INTO app
	(id, appname, url, description, public, sign)
	VALUES
	(0,'auth', '', 'Сервис авторизации', NULL, NULL);

	-- добавление правила, чтобы не нельзя было удалить приложение auth
	CREATE RULE no_delete_app_auth AS ON DELETE TO "app" WHERE id=0 DO INSTEAD NOTHING;

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные app уже существуют.';
END $$;




DO $$
BEGIN

	-- выравнивание автоинкрементного счетчика таблицы app (должно быть в отдельном DO-блоке!)
	SELECT pg_catalog.setval('app_id_seq', COALESCE((SELECT MAX(id)+1 FROM "app"), 1), false);

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Не удалось обновить счетчик app_id_seq.';
END $$;




DO $$
BEGIN

	INSERT INTO "user"
	(id, username    , password                                                          , email              , fullname           , description)
	VALUES
	(0, 'admin'     , '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f', 'admin@rg.ru'      , 'Админ Админов'    , 'Администратор auth-proxy');
	-- пароль для admin: rosgas2011 => '07dd3b6bf9336d7232f7c43fcfcab2c5ae63b7425408c0a7f12b57e638dc6f0f'

	-- добавление правила, чтобы не нельзя было удалить пользователя admin
	CREATE RULE no_delete_user_admin AS ON DELETE TO "user" WHERE id=0 DO INSTEAD NOTHING;

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные user уже существуют.';
END $$;




DO $$
BEGIN

	-- выравнивание автоинкрементного счетчика таблицы user (должно быть в отдельном DO-блоке!)
	SELECT pg_catalog.setval('user_id_seq', COALESCE((SELECT MAX(id)+1 FROM "user"), 1), false);

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Не удалось обновить счетчик user_id_seq.';
END $$;





DO $$
BEGIN

	INSERT INTO app_user_role
	(appname                 , username    , rolename)
	VALUES
	('auth'                  , 'admin'     , 'authadmin');

	-- добавление правила, чтобы не нельзя было удалить роль authadmin, которая разрешает производить изменения в приложении auth
	CREATE RULE no_delete_authadmin AS ON DELETE TO app_user_role WHERE appname='auth' AND username = 'admin' AND rolename = 'authadmin' DO INSTEAD NOTHING;

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные app_user_role уже существуют.';
END $$;





-- Группы ---------------------------------------------------------
DO $$
BEGIN

	INSERT INTO "group"
	(id, groupname       , description)
	VALUES
	(0 ,'admins'         , 'Группа администраторов');



EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные group уже существуют.';
END $$;




DO $$
BEGIN

	-- выравнивание автоинкрементного счетчика таблицы group (должно быть в отдельном DO-блоке!)
	SELECT pg_catalog.setval('group_id_seq', COALESCE((SELECT MAX(id)+1 FROM "group"), 1), false);

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Не удалось обновить счетчик group_id_seq.';
END $$;




DO $$
BEGIN

	INSERT INTO group_user_role
	(group_id, user_id)
	VALUES
	(0, 0); -- связыват пользователя admin и группу admins

EXCEPTION WHEN OTHERS THEN RAISE WARNING 'Данные group_user_role уже существуют.';
END $$;


-- EXCEPTION WHEN OTHERS THEN
--     RAISE EXCEPTION 'Тестовые данные уже существуют.';
-- END$$;
-- COMMIT;
