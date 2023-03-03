

CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS authpublic;
CREATE SCHEMA IF NOT EXISTS extensions;

-- Справочная таблица.
-- Пользователь с дополнительными сведениями
CREATE TABLE IF NOT EXISTS "user" (
    id serial,
    username text NOT NULL,
    password text NOT NULL,
    email text NOT NULL,
    fullname text,
    description text,
    disabled integer NOT NULL DEFAULT 0,

    pashash text,                               -- временный хэш для смены пароля
    pinrequired boolean NOT NULL DEFAULT FALSE, -- требуется ли вводить PIN Google Authenticator для входа в систему

    pinhash_temp text,                          -- Новое значение хэша, которое заменит старое при установке аутентификатора.
                                                -- Наличие ненулевого значения в этом поле сигнализирует:
                                                --   1. установил ли пользователь Google Authenticator на своем телефоне?
                                                --   2. показывать ли ему страницу установки аутентификатора?

	pinhash text,                               -- хэш для первоначальной настройки Google Authenticator
    emailhash text,                             -- хэш для проверки email
    emailconfirmed boolean NOT NULL DEFAULT FALSE, -- подтвержеден ли email пользователя

    CONSTRAINT user_pkey PRIMARY KEY (username),
	CONSTRAINT user_email_unique_idx UNIQUE (email),
	CONSTRAINT user_id_unique_idx UNIQUE (id)
);

CREATE INDEX IF NOT EXISTS user_pashash_idx ON "user" (pashash);

-- Таблица для удаленных пользователей
CREATE TABLE IF NOT EXISTS "user_deleted" (
    id integer,
    username text NOT NULL,
    password text NOT NULL,
    email text NOT NULL,
    fullname text,
    description text,
    disabled integer NOT NULL DEFAULT 0,
    pashash text,                               -- временный хэш для смены пароля
    pinrequired boolean NOT NULL DEFAULT FALSE, -- требуется ли вводить PIN Google Authenticator для входа в систему
    pinhash_temp text,                          -- Новое значение хэша, которое заменит старое при установке аутентификатора.
	pinhash text,                               -- хэш для первоначальной настройки Google Authenticator
    emailhash text,                             -- хэш для проверки email
    emailconfirmed boolean NOT NULL DEFAULT FALSE, -- подтвержеден ли email пользователя
    deleted_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- время удаления
    deleted_by text -- кто удалил
);


-- Справочная таблица.
-- Приложение с описанием
CREATE TABLE IF NOT EXISTS app (
	id serial,
    appname text NOT NULL,
    url text,
    description text,
    rebase text,
    public text,
    sign text,

    CONSTRAINT app_pkey PRIMARY KEY (appname),
	CONSTRAINT app_id_unique_idx UNIQUE (id)
);


----------------------------------------------------------------
-- ОСНОВНАЯ ТАБЛИЦА.
-- Роль пользователя в приложении
CREATE TABLE IF NOT EXISTS app_user_role (
    appname text NOT NULL,
    username text NOT NULL,
    rolename text NOT NULL,

    CONSTRAINT app_user_role_pkey PRIMARY KEY (appname, username, rolename),
    CONSTRAINT app_user_role_fk_u FOREIGN KEY (username) REFERENCES "user"(username) ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE,
    CONSTRAINT app_user_role_fk_a FOREIGN KEY (appname)  REFERENCES "app" (appname)  ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE
);

-- индексы для ускорения выборок

CREATE INDEX IF NOT EXISTS aur_appname_idx ON app_user_role (appname);
CREATE INDEX IF NOT EXISTS aur_username_idx ON app_user_role (username);



-- GROUPSGROUPSGROUPSGROUPSGROUPSGROUPSGROUPSGROUPSGROUPSGROUPSGROUPSGROUPS


-- Группа
CREATE TABLE IF NOT EXISTS "group" (
    id serial,
    groupname text NOT NULL,
    description text,
    disabled integer NOT NULL DEFAULT 0,

    CONSTRAINT group_pkey PRIMARY KEY (id),
    UNIQUE (groupname)
);



-- Роль пользователя в группе
CREATE TABLE IF NOT EXISTS group_user_role (
    group_id int NOT NULL,
    user_id int NOT NULL,
    rolename text NOT NULL DEFAULT 'member',

    CONSTRAINT group_user_role_pkey PRIMARY KEY (group_id, user_id, rolename),
    CONSTRAINT group_user_role_fk_u FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE,
    CONSTRAINT group_user_role_fk_g FOREIGN KEY (group_id)  REFERENCES "group" (id)  ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE
);


-- Роль группы в приложении
CREATE TABLE IF NOT EXISTS group_app_role (
    group_id int NOT NULL,
    app_id int NOT NULL,
    rolename text NOT NULL,

    CONSTRAINT group_app_role_pkey PRIMARY KEY (group_id, app_id, rolename),
    CONSTRAINT group_app_role_fk_g FOREIGN KEY (group_id) REFERENCES "group"(id) ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE,
    CONSTRAINT group_app_role_fk_a FOREIGN KEY (app_id)  REFERENCES "app" (id)  ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE
);


