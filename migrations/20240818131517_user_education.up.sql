create type user_education_type_enum as enum ('SCHOOL', 'COLLAGE');

create table
  public.user_educations (
    id bigserial,
    user_id uuid not null,
    order_index smallint not null,
    type user_education_type_enum not null,
    institute_name text not null,
    grade decimal not null,
    attributes jsonb,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint user_education_pkey primary key (id),
    constraint user_education_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade,
    constraint user_education_user_id_and_order_index_composite_key unique (user_id, order_index)
  ) tablespace pg_default;