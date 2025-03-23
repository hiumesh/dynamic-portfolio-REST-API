create table
  public.hackathons (
    id bigserial,
    user_id uuid not null,
    order_index smallint not null,
    avatar text,
    title text not null,
    location text not null,
    start_date date not null,
    end_date date not null,
    description text not null,
    certificate_link text,
    attributes jsonb not null default '{}',
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint hackathons_pkey primary key (id),
    constraint hackathons_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade,
    constraint hackathons_user_id_and_order_index_composite_key unique (user_id, order_index) deferrable initially deferred
  ) tablespace pg_default;