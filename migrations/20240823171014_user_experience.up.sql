create type user_experiences_job_type_enum as enum ('PART_TIME', 'SEMI_FULL_TIME', 'FULL_TIME');


create table
  public.user_experiences (
    id bigserial,
    user_id uuid not null,
    order_index smallint not null,
    company_name text not null,
    company_url text,
    job_type user_experiences_job_type_enum not null,
    job_title text not null,
    location text not null,
    start_date date not null,
    end_date date,
    description text[] not null,
    skills_used text[] not null,
    certificate_link text,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint user_experiences_pkey primary key (id),
    constraint user_experiences_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade,
    constraint user_experiences_user_id_and_order_index_composite_key unique (user_id, order_index) deferrable initially deferred
  ) tablespace pg_default;