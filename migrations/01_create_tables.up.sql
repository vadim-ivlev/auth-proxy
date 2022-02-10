

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

    CONSTRAINT user_pkey PRIMARY KEY (username)
);

-- Патч для таблицы пользователей чтобы обеспечить вход по PIN,
-- И более безопасную смену пароля.
ALTER TABLE "user"
    -- временный хэш для смены пароля
    ADD COLUMN IF NOT EXISTS pashash text, 
    -- требуется ли вводить PIN Google Authenticator для входа в систему
    ADD COLUMN IF NOT EXISTS pinrequired boolean NOT NULL DEFAULT FALSE,  
    -- Новое значение хэша, которое заменит старое при установке аутентификатора. 
    -- Наличие ненулевого значения в этом поле сигнализирует:
    --   1. установил ли пользователь Google Authenticator на своем телефоне?
    --   2. показывать ли ему страницу установки аутентификатора?
    ADD COLUMN IF NOT EXISTS pinhash_temp text,
    -- хэш для первоначальной настройки Google Authenticator 
    ADD COLUMN IF NOT EXISTS pinhash text, -- DEFAULT uuid_in(md5(random()::text || clock_timestamp()::text)::cstring)::text ;
    -- хэш для проверки email
    ADD COLUMN IF NOT EXISTS emailhash text,
    -- подтвержеден ли email пользователя
    ADD COLUMN IF NOT EXISTS emailconfirmed boolean NOT NULL DEFAULT FALSE;



-- Справочная таблица.
-- Приложение с описанием
CREATE TABLE IF NOT EXISTS app (
    appname text NOT NULL,
    url text,
    description text,
    rebase text,
    public text,
    sign text,

    CONSTRAINT app_pkey PRIMARY KEY (appname)
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
CREATE INDEX IF NOT EXISTS user_pashash_idx ON "user" (pashash);
CREATE INDEX IF NOT EXISTS aur_appname_idx ON app_user_role (appname);
CREATE INDEX IF NOT EXISTS aur_username_idx ON app_user_role (username);
CREATE UNIQUE INDEX IF NOT EXISTS user_email_unique_idx ON "user" (email);
CREATE UNIQUE INDEX IF NOT EXISTS user_id_unique_idx ON "user" (id);
-- CREATE INDEX IF NOT EXISTS user_textsearch_idx ON "user" USING gin (to_tsvector('russian', fullname || ' ' || description  || ' ' || email  || ' ' || username ));
-- CREATE INDEX IF NOT EXISTS app_textsearch_idx ON "app" USING gin (to_tsvector('russian', appname || ' ' || description ));



-- --------------------------------------------------------------

ALTER TABLE "app" ADD COLUMN IF NOT EXISTS id serial;
-- ALTER TABLE "app" ADD UNIQUE (id);


ALTER TABLE "app" DROP CONSTRAINT IF EXISTS app_unique_id_constraint CASCADE;
ALTER TABLE "app" ADD CONSTRAINT app_unique_id_constraint UNIQUE (id);



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
    rolename text NOT NULL,

    CONSTRAINT group_user_role_pkey PRIMARY KEY (group_id, user_id, rolename),
    CONSTRAINT group_user_role_fk_u FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE,
    CONSTRAINT group_user_role_fk_g FOREIGN KEY (group_id)  REFERENCES "group" (id)  ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE
);


-- Роль группы в приложении
CREATE TABLE IF NOT EXISTS group_app_role (
    id serial NOT NULL,
    group_id int NOT NULL,
    app_id int NOT NULL,
    rolename text NOT NULL

    ,CONSTRAINT group_app_role_pkey PRIMARY KEY (id)
    ,UNIQUE (group_id, app_id, rolename)
    -- ,CONSTRAINT group_app_role_pkey PRIMARY KEY (group_id, app_id, rolename)
    ,CONSTRAINT group_app_role_fk_g FOREIGN KEY (group_id) REFERENCES "group"(id) ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE
    ,CONSTRAINT group_app_role_fk_a FOREIGN KEY (app_id)  REFERENCES "app" (id)  ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE
);


