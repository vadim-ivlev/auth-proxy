

CREATE SCHEMA IF NOT EXISTS auth; 
CREATE SCHEMA IF NOT EXISTS extensions; 

-- Справочная таблица.
-- Пользователь с дополнительными сведениями 
CREATE TABLE IF NOT EXISTS "user" (
    username text NOT NULL,
    password text NOT NULL,
    email text NOT NULL,
    fullname text,
    description text,
    disabled integer NOT NULL DEFAULT 0,
    id serial,
    -- требуется ли вводить PIN Google Authenticator для входа в систему
    pinrequired boolean NOT NULL DEFAULT FALSE,  
    -- установил ли пользователь Google Authenticator на своем телефоне.
    -- показывать ли ему страницу установки аутентификатора.
    pinset boolean NOT NULL DEFAULT FALSE,  
    pinhash text, -- хэш для первоначальной настройки Google Authenticator

    CONSTRAINT user_pkey PRIMARY KEY (username)
);

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
CREATE INDEX IF NOT EXISTS aur_appname_idx ON app_user_role (appname);
CREATE INDEX IF NOT EXISTS aur_username_idx ON app_user_role (username);
CREATE UNIQUE INDEX IF NOT EXISTS user_email_unique_idx ON "user" (email);
CREATE UNIQUE INDEX IF NOT EXISTS user_id_unique_idx ON "user" (id);
-- CREATE INDEX IF NOT EXISTS user_textsearch_idx ON "user" USING gin (to_tsvector('russian', fullname || ' ' || description  || ' ' || email  || ' ' || username ));
-- CREATE INDEX IF NOT EXISTS app_textsearch_idx ON "app" USING gin (to_tsvector('russian', appname || ' ' || description ));


