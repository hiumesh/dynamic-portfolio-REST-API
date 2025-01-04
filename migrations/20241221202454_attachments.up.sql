create table
  public.attachments (
    id bigserial,
    -- order_index smallint not null,
    parent_table text not null,
    parent_id bigint not null,
    user_id uuid not null,
    file_url text not null,
    file_name text not null,
    file_type text not null,
    file_size bigint not null,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint attachments_pkey primary key (id),
    constraint attachments_file_url_ukey unique (file_url),
    constraint attachments_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade
    -- constraint attachments_parent_id_and_order_index_and_parent_table_composite_key unique (parent_id, order_index, parent_table) deferrable initially deferred
  ) tablespace pg_default;