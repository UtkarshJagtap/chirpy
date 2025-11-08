-- +goose Up
CREATE TABLE refresh_token (
  token text PRIMARY KEY,
  created_at timestamp not null,
  updated_at timestamp not null,
  user_id uuid not null,
  expires_at timestamp not null,
  revoked_at timestamp,
  FOREIGN KEY(user_id) REFERENCES users(id) on delete cascade
);

-- +goose Down
DROP TABLE refresh_token;
