-- TRUNCATE public.broadcast CASCADE;

DELETE FROM public.broadcast CASCADE WHERE id < 360;
