create table
  public.certifications (
    id bigserial,
    user_id uuid not null,
    order_index smallint not null,
    title text not null,
    description text not null,
    completion_date date,
    certificate_link text,
    skills_used text[] not null,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint certifications_pkey primary key (id),
    constraint certifications_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade,
    constraint certifications_user_id_and_order_index_composite_key unique (user_id, order_index) deferrable initially deferred
  ) tablespace pg_default;