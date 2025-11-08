-- +goose Up
ALTER TABLE users
  ADD COLUMN chirpy_red bool not null default false;

-- +goose Down
ALTER TABLE users
  DROP COLUMN chirpy_red;


