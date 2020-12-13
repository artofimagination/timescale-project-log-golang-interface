-- +migrate Up
CREATE TABLE IF NOT EXISTS projects_data(
   created_at timestamp NOT NULL DEFAULT NOW() PRIMARY KEY,
   viewer_id int NOT NULL,
   data json
);

-- +migrate Up
CREATE INDEX ON projects_data (created_at DESC, viewer_id);