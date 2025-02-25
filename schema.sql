create table users (
  id text not null primary key,
  username text not null unique,
  profile_picture text not null,
  created_at timestamp not null,
  updated_at timestamp not null
);
