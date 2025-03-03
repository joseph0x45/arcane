create table users (
  id text not null primary key,
  github_id text not null unique,
  username text not null,
  avatar_url text not null,
  joined_at timestamp not null
);

create table sessions (
  id text not null primary key,
  user_id text not null references users(id),
  is_valid bool not null
);
