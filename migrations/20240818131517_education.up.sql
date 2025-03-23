create type educations_type_enum as enum ('SCHOOL', 'COLLEGE');

create table
  public.educations (
    id bigserial,
    user_id uuid not null,
    order_index smallint not null,
    type educations_type_enum not null,
    institute_name text not null,
    grade decimal not null,
    attributes jsonb,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint educations_pkey primary key (id),
    constraint educations_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade,
    constraint educations_user_id_and_order_index_composite_key unique (user_id, order_index) deferrable initially deferred
  ) tablespace pg_default;