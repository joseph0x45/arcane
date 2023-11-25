create table users (
  id uuid not null primary key,
  github_id text not null unique,
  email text not null unique,
  username text not null,
  avatar_url text not null
);

create table teams (
  id uuid not null primary key,
  name text not null,
  owner uuid not null references users(id),
  plan text not null default 'basic'
);

create table team_memberships(
  id uuid not null primary key,
  team uuid not null references teams(id),
  member uuid not null references users(id),
  position text not null default 'member'
);
