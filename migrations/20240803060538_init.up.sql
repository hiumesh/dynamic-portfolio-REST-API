-- types

create type user_profiles_portfolio_status_enum as enum ('DRAFT', 'ACTIVE', 'IN_ACTIVE');

-- tables

create table public.user_profiles (
  user_id uuid not null,
  email text not null,
  full_name text null,
  avatar_url text null,
  slug text null,
  portfolio_status public.user_profiles_portfolio_status_enum not null default 'DRAFT'::user_profiles_portfolio_status_enum,
  attributes jsonb not null default '{}'::jsonb,
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,
  deleted_at timestamp with time zone null,
  fts tsvector GENERATED ALWAYS as (
    to_tsvector(
      'english'::regconfig,
      (
        (
          (
            (COALESCE(full_name, ''::text) || ' '::text) || email
          ) || ' '::text
        ) || COALESCE((attributes ->> 'college'::text), ''::text)
      )
    )
  ) STORED null,
  constraint user_profiles_pkey primary key (user_id),
  constraint user_profiles_email_ukey unique (email),
  constraint user_profiles_slug_ukey unique (slug),
  constraint user_profiles_user_id_fkey foreign KEY (user_id) references auth.users (id) on delete CASCADE
) tablespace pg_default;

-- indexes

create index user_profiles_fts_idx on user_profiles using gin (fts);

-- functions

create or replace function insert_user_profiles()
returns trigger as $$
begin
    insert into public.user_profiles (user_id, email, created_at, updated_at)
    values (NEW.id, NEW.email, now(), now());
    
    return NEW;
end;
$$ language plpgsql security definer;

-- triggers

create trigger after_insert_users
after insert on auth.users
for each row
execute function insert_user_profiles();