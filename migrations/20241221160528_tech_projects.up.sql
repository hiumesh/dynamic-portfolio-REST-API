create table
  public.tech_projects (
    id bigserial,
    user_id uuid not null,
    order_index smallint not null,
    title text not null,
    description text not null,
    tech_used text[] not null,
    attributes jsonb not null default '{}',
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint tech_projects_pkey primary key (id),
    constraint tech_projects_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade,
    constraint tech_projects_user_id_and_order_index_composite_key unique (user_id, order_index) deferrable initially deferred
  ) tablespace pg_default;

  create index tech_projects_fts_idx on tech_projects using gin (fts);
