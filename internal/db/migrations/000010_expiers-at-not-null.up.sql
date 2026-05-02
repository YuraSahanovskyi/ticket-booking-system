UPDATE public.bookings 
SET expires_at = created_at + INTERVAL '24 hours' 
WHERE expires_at IS NULL;

ALTER TABLE public.bookings ALTER COLUMN expires_at SET NOT NULL;