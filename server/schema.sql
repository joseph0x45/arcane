create table users (
  id text not null primary key,
  email text not null unique,
  password text not null
);

create table sessions (
  id text not null primary key,
  user_id text not null references users(id),
  is_valid bool not null
);
