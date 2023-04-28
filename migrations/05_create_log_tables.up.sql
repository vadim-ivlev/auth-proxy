

-- логи запросов на создание пользователей
CREATE TABLE IF NOT EXISTS "create_user_log" (
    id serial NOT NULL PRIMARY KEY,
    email text,
    password text,
    fullname text,
    description text,
    ip text,
    full_path text,  -- полный путь запроса
    user_agent text, -- user agent браузера
    referer text,    -- с какой страницы пришел запрос
    headers jsonb,   -- заголовки запроса, включая куки. SELECT  headers->'Connection' ... FROM create_user_log;
    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);



