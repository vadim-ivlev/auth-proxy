
-- Справочная таблица.
-- Пользователь с дополнительными сведениями 
CREATE TABLE IF NOT EXISTS "user" (
    username text NOT NULL,
    password text NOT NULL,
    email text,
    fullname text,
    description text,

    CONSTRAINT user_pkey PRIMARY KEY (username)
);

-- Справочная таблица.
-- Приложение с описанием
CREATE TABLE IF NOT EXISTS app (
    appname text NOT NULL,
    description text,

    CONSTRAINT app_pkey PRIMARY KEY (appname)
);

-- Справочная таблица.
-- Роль приложения с описанием
CREATE TABLE IF NOT EXISTS app_role (
    appname text NOT NULL,
    rolename text NOT NULL,
    description text,

    CONSTRAINT role_pkey PRIMARY KEY (appname, rolename),
    CONSTRAINT role_fk FOREIGN KEY (appname) REFERENCES app(appname) ON DELETE CASCADE DEFERRABLE
);

----------------------------------------------------------------
-- ОСНОВНАЯ ТАБЛИЦА. 
-- Роль пользователя в приложении
CREATE TABLE IF NOT EXISTS app_user_role (
    appname text NOT NULL,
    username text NOT NULL,
    rolename text NOT NULL,

    CONSTRAINT app_user_role_pkey PRIMARY KEY (appname, username, rolename),
    CONSTRAINT app_user_role_fk_u FOREIGN KEY (username) REFERENCES "user"(username) ON DELETE CASCADE DEFERRABLE,
    CONSTRAINT app_user_role_fk_ar FOREIGN KEY (appname, rolename) REFERENCES app_role(appname, rolename) ON DELETE CASCADE DEFERRABLE
);

-- индексы для ускорения выборок
CREATE INDEX IF NOT EXISTS role_appname_idx ON app_role USING btree (appname);
CREATE INDEX IF NOT EXISTS user_textsearch_idx ON "user" USING gin (to_tsvector('russian', fullname || ' ' || description  || ' ' || email  || ' ' || username ));


