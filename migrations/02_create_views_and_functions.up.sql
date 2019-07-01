-- ФУНКЦИИ **********************************************************************************
-- Назначение функций преобразовть таблицы в JSON требуемый API RG.RU.
-- Выражение array_to_json(array_agg(row_to_json( t, false )),true) преобразует таблицу t JSON объект.
-- Выражение row_to_json( t, false ) преобразует одну запись таблицы t JSON объект.
-- Все функции возвращают JSON.


-- NOTE: 
-- Calls jsonb_agg(t) seem to produce the same results as   
-- array_to_json(array_agg(row_to_json( t, false )),true),
-- which were used in onlinebc_admin.



-- R E S T functions  ************************************************************************

CREATE OR REPLACE FUNCTION public.get_images(idd integer)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN
    (
        select jsonb_agg(t) from
        (
            SELECT * FROM image
            WHERE post_id = idd
        ) t
    );

END;
$function$
;


CREATE OR REPLACE FUNCTION public.get_answers(idd integer)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN
    (
        select jsonb_agg(t) from
        (
            SELECT * FROM answer_view
            WHERE id_parent = idd
        ) t
    );

END;
$function$
;

CREATE OR REPLACE FUNCTION public.get_posts(idd integer)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
BEGIN   
    RETURN  
    (
        select jsonb_agg(t) from
        ( 
            SELECT * FROM post_view
            WHERE 
                    id_broadcast = idd 
                AND id_parent IS NULL 
            ORDER BY id DESC 
        ) t 
    );
    
END;
$function$
;


-- get_full_broadcast - гланая функция.
-- Преобразует плоские таблицы базы данных в JSON
-- с многоуровненвой иерархической структурой

CREATE OR REPLACE FUNCTION public.get_full_broadcast(idd integer)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
BEGIN   
    RETURN
    (
        select jsonb_agg(t) from
        ( select *, get_posts(id) as posts  from broadcast where id = idd ) t
    );
END;
$function$
;

CREATE OR REPLACE FUNCTION public.get_broadcasts()
 RETURNS json
 LANGUAGE plpgsql
AS $function$
BEGIN   
    RETURN
    (
       select jsonb_agg(t) from
       ( select *  from broadcast ) t
    );
END;
$function$
;


CREATE OR REPLACE FUNCTION public.get_broadcast_list(where_part text default '')
 RETURNS json
 LANGUAGE plpgsql
AS $function$
DECLARE
    j json;
BEGIN   
    EXECUTE '( 
                select jsonb_agg(t) FROM 
                ( SELECT * FROM broadcast '|| where_part ||' ) t 
             )' INTO j;
    RETURN j;
END;
$function$
;




-- R E S T ПРЕДСТАВЛЕНИЯ *****************************************************************************888
-- именуют поля таблицы post в соответствии с полями JSON API RG.RU
-- https://outer.rg.ru/plain/online_translations/api/online.php?id=354
-- Если нужно поменять названия полей, добавить или убрать поле
-- менять нужно здесь.


CREATE OR REPLACE VIEW public.answer_view
AS SELECT post.id,
    post.id_parent,
    post.id_broadcast,
    post.has_big_img,
    post.author     AS posts__answer__author,
    post.text       AS posts__answer__text,
    post.post_type  AS posts__answer__type,
    post.link       AS posts__answer__uri,
    post.post_time  AS posts__answer__timestamp,
    -- to_char(to_timestamp(post.post_time::double precision), 'DD.MM.YYYY'::text) AS posts__answer__date,
    -- to_char(to_timestamp(post.post_time::double precision), 'HH24:MI'::text)    AS posts__answer__time,
    
    get_images(post.id)      AS posts__answer__images,
    get_answers(post.id)    AS posts__answer__answers
   FROM post;



CREATE OR REPLACE VIEW public.post_view
AS SELECT post.id,
    post.id_parent,
    post.id_broadcast,
    post.has_big_img,
    post.author AS posts__author,
    post.text AS posts__text,
    post.post_type AS posts__type,
    post.link AS posts__uri,
    post.post_time AS posts__timestamp,
    -- to_char(to_timestamp(post.post_time::double precision), 'DD.MM.YYYY'::text) AS posts__date,
    -- to_char(to_timestamp(post.post_time::double precision), 'HH24:MI'::text) AS posts__time,
    
    get_images(post.id) AS posts__images,
    get_answers(post.id) AS posts__answers
   FROM post;




-- G R A P H Q L 

-- G R A P H Q L functions  ************************************************************************

CREATE OR REPLACE FUNCTION public.get_full_post_images(idd integer)
    RETURNS json
    LANGUAGE plpgsql
    AS $function$
    BEGIN
        RETURN
        (
            select jsonb_agg(t) from
            (
                SELECT * FROM image
                WHERE post_id = idd
            ) t
        );

    END;
    $function$
;

CREATE OR REPLACE FUNCTION public.get_full_post_answers(idd integer)
    RETURNS json
    LANGUAGE plpgsql
    AS $function$
    BEGIN
        RETURN
        (
            select jsonb_agg(t) from
            (
                SELECT * FROM full_answer
                WHERE id_parent = idd
            ) t
        );

    END;
    $function$
;

CREATE OR REPLACE FUNCTION public.get_full_broadcast_posts(idd integer)
    RETURNS json
    LANGUAGE plpgsql
    AS $function$
    BEGIN   
        RETURN  
        (
            select jsonb_agg(t) from
            ( 
                SELECT * FROM full_post
                WHERE 
                        id_broadcast = idd 
                    AND id_parent IS NULL 
                ORDER BY id DESC 
            ) t 
        );
        
    END;
    $function$
;




-- G R A P H Q L  V I E W S   *****************************************************************************888

CREATE OR REPLACE VIEW public.full_answer AS
    SELECT  * 
            , get_full_post_images(post.id) AS images 
    FROM post 
;

CREATE OR REPLACE VIEW public.full_post AS 
    SELECT  *
            , get_full_post_images(post.id) AS images
            , get_full_post_answers(post.id) AS answers
    FROM post
;

CREATE OR REPLACE VIEW public.full_broadcast AS 
    SELECT  *
            , get_full_broadcast_posts(broadcast.id) AS posts
    FROM broadcast
;

