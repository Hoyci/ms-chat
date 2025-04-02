CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS contacts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id uuid NOT NULL,
    contact_id uuid NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'pending',
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp DEFAULT null,
    deleted_at timestamp DEFAULT null
);