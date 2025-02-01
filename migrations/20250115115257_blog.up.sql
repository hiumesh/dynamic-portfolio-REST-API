create table
  public.blogs (
    id bigserial,
    user_id uuid not null,
    cover_image text,
    title text not null,
    body text,
    slug text not null,
    attributes jsonb not null default '{}',
    published_at timestamptz,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint blogs_pkey primary key (id),
    constraint blogs_slug_ukey unique (slug),
    constraint blogs_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade
  ) tablespace pg_default;

create table public.tags (
    id bigserial,
    name text not null,
    description text,
    user_id uuid,
    attributes jsonb not null default '{}',
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint tags_pkey primary key (id),
    constraint tags_name_ukey unique (name),
    constraint tags_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade
) tablespace pg_default;

create table public.blog_tags (
    blog_id bigint not null,
    tag_id bigint not null,
    constraint blog_tags_pkey primary key (blog_id, tag_id),
    constraint blog_tags_blog_id_fkey foreign key (blog_id) references public.blogs (id) on delete cascade,
    constraint blog_tags_tag_id_fkey foreign key (tag_id) references public.tags (id) on delete cascade
) tablespace pg_default;