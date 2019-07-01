CREATE TABLE IF NOT EXISTS public.broadcast (
    id integer NOT NULL GENERATED BY DEFAULT AS IDENTITY (START WITH 1000),
    title character varying(256),
    time_created integer,
    time_begin integer,
    is_ended integer,
    show_date integer,
    show_time integer,
    show_main_page integer,
    link_article character varying(256),
    link_img character varying(255),
    groups_create integer,
    is_diary integer,
    diary_author character varying(255),

    CONSTRAINT broadcast_pkey1 PRIMARY KEY (id)
);



CREATE TABLE IF NOT EXISTS public.post (
    id integer NOT NULL GENERATED BY DEFAULT AS IDENTITY (START WITH 30000),
    id_parent integer,
    id_broadcast integer,
    text text,
    post_time integer,
    post_type integer,
    link character varying(256),
    has_big_img integer,
    author text,

    CONSTRAINT post_pkey1 PRIMARY KEY (id),
    CONSTRAINT post_broadcast_fk FOREIGN KEY (id_broadcast) REFERENCES public.broadcast(id) ON DELETE CASCADE DEFERRABLE,
    CONSTRAINT post_post_fk FOREIGN KEY (id_parent) REFERENCES public.post(id) ON DELETE CASCADE DEFERRABLE
);



CREATE TABLE IF NOT EXISTS public.image (
    id integer NOT NULL GENERATED BY DEFAULT AS IDENTITY (START WITH 6000),
    post_id integer NOT NULL,
    filepath character varying(255),
    thumbs jsonb,
    source character varying(255),
    width integer,
    height integer,

    CONSTRAINT image_pk PRIMARY KEY (id),
    CONSTRAINT image_post_fk FOREIGN KEY (post_id) REFERENCES public.post(id) ON DELETE CASCADE DEFERRABLE
);


CREATE INDEX IF NOT EXISTS image_post_id_idx ON public.image USING btree (post_id);
CREATE INDEX IF NOT EXISTS post_id_broadcast_idx ON public.post USING btree (id_broadcast);
CREATE INDEX IF NOT EXISTS post_id_parent_idx ON public.post USING btree (id_parent);
CREATE INDEX IF NOT EXISTS broadcast_title_textsearch_idx ON public.broadcast USING gin (to_tsvector('russian', title));

