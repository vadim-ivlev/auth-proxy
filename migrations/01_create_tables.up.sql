
-- Пользователи  
CREATE TABLE IF NOT EXISTS public.user (
    username text NOT NULL,
    password text NOT NULL,
    email text,
    fullname text,
    description text,

    CONSTRAINT user_pkey PRIMARY KEY (username)
);

-- Приложения к которым требуется авторизация
CREATE TABLE IF NOT EXISTS public.app (
    appname text NOT NULL,
    description text,
    url text,

    CONSTRAINT app_pkey PRIMARY KEY (appname)
);

-- Роли приложения
CREATE TABLE IF NOT EXISTS public.app_role (
    appname text NOT NULL,
    rolename text NOT NULL,
    description text,

    CONSTRAINT role_pkey PRIMARY KEY (appname, rolename),
    CONSTRAINT role_fk FOREIGN KEY (appname) REFERENCES public.app(appname) ON DELETE CASCADE DEFERRABLE
);


-- Роли пользователя для приложения
CREATE TABLE IF NOT EXISTS public.app_user_role (
    appname text NOT NULL,
    username text NOT NULL,
    rolename text NOT NULL,

    CONSTRAINT app_user_role_pkey PRIMARY KEY (appname, username, rolename),
    CONSTRAINT app_user_role_fk_u FOREIGN KEY (username) REFERENCES public.user(username) ON DELETE CASCADE DEFERRABLE,
    CONSTRAINT app_user_role_fk_ar FOREIGN KEY (appname, rolename) REFERENCES public.app_role(appname, rolename) ON DELETE CASCADE DEFERRABLE
);

CREATE INDEX IF NOT EXISTS role_appname_idx ON public.app_role USING btree (appname);
CREATE INDEX IF NOT EXISTS user_textsearch_idx ON public.user USING gin (to_tsvector('russian', fullname || ' ' || description  || ' ' || email  || ' ' || username ));


