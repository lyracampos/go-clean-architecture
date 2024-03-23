CREATE TABLE users (
      id bigserial PRIMARY KEY,
      first_name text NOT NULL,
      last_name text NOT NULL,
      email text NOT NULL,
      role text NOT NULL,
      created_at timestamp NOT NULL DEFAULT now(),
      updated_at timestamp NOT NULL DEFAULT now()
);