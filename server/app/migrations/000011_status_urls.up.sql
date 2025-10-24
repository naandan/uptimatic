ALTER TABLE urls
ADD COLUMN public_id uuid NOT NULL DEFAULT gen_random_uuid();