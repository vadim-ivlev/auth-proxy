-- GRAPHQL ---------------------------------
DROP VIEW IF EXISTS public.full_broadcast;
DROP VIEW IF EXISTS public.full_post;
DROP VIEW IF EXISTS public.full_answer;

DROP FUNCTION IF EXISTS public.get_full_broadcast_posts(idd integer);
DROP FUNCTION IF EXISTS public.get_full_post_answers(idd integer);
DROP FUNCTION IF EXISTS public.get_full_post_images(idd integer);




-- REST ------------------------------------
DROP VIEW IF EXISTS public.post_view;
DROP VIEW IF EXISTS public.answer_view;

DROP FUNCTION IF EXISTS public.get_broadcast_list(where_part text);
DROP FUNCTION IF EXISTS public.get_broadcasts();
DROP FUNCTION IF EXISTS public.get_full_broadcast(idd integer);
DROP FUNCTION IF EXISTS public.get_posts(idd integer);
DROP FUNCTION IF EXISTS public.get_answers(idd integer);
DROP FUNCTION IF EXISTS public.get_images(idd integer);

