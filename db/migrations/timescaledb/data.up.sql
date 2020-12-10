-- +migrate Up
CREATE extension IF NOT EXISTS "uuid-ossp";

-- +migrate Up
CREATE TABLE IF NOT EXISTS projects(
   id uuid NOT NULL,
   users_id uuid NOT NULL,
   run_seq_no integer NOT NULL DEFAULT 0,
   data json
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS projects_data(
   created_at timestamp NOT NULL DEFAULT NOW() PRIMARY KEY,
   projects_id uuid NOT NULL,
   FOREIGN KEY (projects_id) REFERENCES projects(id),
   run_seq_no integer NOT NULL DEFAULT 0,
   data json
);

-- +migrate Up
CREATE INDEX ON projects_data (created_at DESC, project_id);