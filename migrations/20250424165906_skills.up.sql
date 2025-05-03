create table public.skills (
    id bigserial,
    name text not null,
    image text,
    constraint skills_pkey primary key (id),
    constraint skills_name_ukey unique (name)
) tablespace pg_default;