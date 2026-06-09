-- Schema for the Account context.
-- Run the trigger block once in the Supabase SQL Editor to auto-create
-- accounts on new Supabase Auth user registration.

CREATE TABLE IF NOT EXISTS public.accounts (
    id              TEXT PRIMARY KEY,  -- Supabase Auth user_id
    whatsapp_number TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Trigger: creates an Account row for every new Supabase Auth user.
-- Run once manually in Supabase SQL Editor.
--
-- CREATE OR REPLACE FUNCTION public.handle_new_user()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     INSERT INTO public.accounts (id, created_at, updated_at)
--     VALUES (NEW.id, now(), now());
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql SECURITY DEFINER;
--
-- CREATE TRIGGER on_auth_user_created
--     AFTER INSERT ON auth.users
--     FOR EACH ROW EXECUTE FUNCTION public.handle_new_user();
