CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE events ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE seats ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE bookings ALTER COLUMN id SET DEFAULT uuid_generate_v4();