-- +migrate Up
CREATE extension IF NOT EXISTS "uuid-ossp";

-- +migrate Up
CREATE TABLE IF NOT EXISTS projects_data(
   created_at timestamp NOT NULL DEFAULT NOW() PRIMARY KEY,
   viewer_id uuid NOT NULL,
   data jsonb
);

-- +migrate Up
CREATE INDEX ON projects_data (created_at DESC, viewer_id);