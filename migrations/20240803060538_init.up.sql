-- types

create type user_profiles_portfolio_status_enum as enum ('DRAFT', 'ACTIVE', 'IN_ACTIVE');

-- tables

create table
  public.user_profiles (
    user_id uuid not null,
    email text not null,
    full_name text,
    avatar_url text,
    slug text,
    portfolio_status user_profiles_portfolio_status_enum not null default 'DRAFT',
    attributes jsonb not null default '{}',
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint user_profiles_pkey primary key (user_id),
    constraint user_profiles_email_ukey unique (email),
    constraint user_profiles_slug_ukey unique (slug),
    constraint user_profiles_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade
  ) tablespace pg_default;

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