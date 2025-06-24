create type user_reaction_type_enum as enum ('like', 'clap', 'heart');

create table public.comments (
    id bigserial,
    user_id uuid not null,
    parent_id bigint,
    body text not null,
    attributes jsonb not null default '{}',
    created_at timestamptz not null,
    updated_at timestamptz not null,
    deleted_at timestamptz,
    constraint comments_pkey primary key (id),
    constraint comments_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade,
    constraint comments_parent_id_fkey foreign key (parent_id) references public.comments (id) on delete cascade
) tablespace pg_default;

create table public.comment_reactions (
    id bigserial,
    user_id uuid not null,
    comment_id bigint not null,
    type user_reaction_type_enum not null,
    constraint comment_reactions_pkey primary key (id),
    constraint comment_reactions_comment_id_fkey foreign key (comment_id) references public.comments (id) on delete cascade,
    constraint comment_reactions_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade,
    constraint comment_reactions_user_id_and_comment_id_and_type_composite_key unique (user_id, comment_id, type) deferrable initially deferred
) tablespace pg_default;


create table public.blog_comments (
  blog_id bigint not null,
  comment_id bigint not null,
  constraint blog_comments_pkey primary key (blog_id, comment_id),
  constraint blog_comments_blog_id_fkey foreign key (blog_id) references public.blogs (id) on delete cascade,
  constraint blog_comments_comment_id_fkey foreign key (comment_id) references public.comments (id) on delete cascade
) tablespace pg_default;

create table public.blog_reactions (
    id bigserial,
    user_id uuid not null,
    blog_id bigint not null,
    type user_reaction_type_enum not null,
    constraint blog_reactions_pkey primary key (id),
    constraint blog_reactions_blog_id_fkey foreign key (blog_id) references public.blogs (id) on delete cascade,
    constraint blog_reactions_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade,
    constraint blog_reactions_user_id_and_blog_id_and_type_composite_key unique (user_id, blog_id, type) deferrable initially deferred
) tablespace pg_default;

create or replace function update_blog_comments_count()
returns trigger as $$
declare
    comments_count integer;
    blog_id_to_use integer;
begin

    blog_id_to_use := coalesce(NEW.blog_id, OLD.blog_id);

    select count(*) into comments_count
    from public.blog_comments
    where blog_id = blog_id_to_use;

    update public.blogs
    set attributes = jsonb_set(
        attributes,
        '{comments_count}',
        to_jsonb(comments_count),
        true
    )
    where id = blog_id_to_use;

    return null;
end;
$$ language plpgsql;


create trigger update_blog_comments_count_trigger
after insert or delete or update on public.blog_comments
for each row
execute function update_blog_comments_count();

create or replace function update_blog_reactions_count()
returns trigger as $$
declare
    reaction_data jsonb;
    blog_id_to_use integer;
begin
    blog_id_to_use := coalesce(NEW.blog_id, OLD.blog_id);

  
    select jsonb_object_agg(type, cnt) into reaction_data
    from (
        select type, count(*) as cnt
        from public.blog_reactions
        where blog_id = blog_id_to_use
        group by type
    ) as sub;

  
    update public.blogs
    set attributes = jsonb_set(
        attributes,
        '{reaction_metadata}',
        coalesce(reaction_data, '{}'::jsonb),
        true
    )
    where id = blog_id_to_use;

    return null;
end;
$$ language plpgsql;


create trigger update_blog_reactions_count_trigger
after insert or delete or update on public.blog_reactions
for each row
execute function update_blog_reactions_count();

create or replace function update_comment_reactions_count()
returns trigger as $$
declare
    reaction_data jsonb;
    comment_id_to_use integer;
begin
    comment_id_to_use := coalesce(NEW.comment_id, OLD.comment_id);

    select jsonb_object_agg(type, cnt) into reaction_data
    from (
        select type, count(*) as cnt
        from public.comment_reactions
        where comment_id = comment_id_to_use
        group by type
    ) as sub;

    update public.comments
    set attributes = jsonb_set(
        attributes,
        '{reaction_metadata}',
        coalesce(reaction_data, '{}'::jsonb),
        true
    )
    where id = comment_id_to_use;

    return null;
end;
$$ language plpgsql;

create trigger update_comment_reactions_count_trigger
after insert or delete or update on public.comment_reactions
for each row
execute function update_comment_reactions_count();

create or replace function update_comment_replies_count()
returns trigger as $$
declare
    reply_count integer;
    parent_comment_id integer;
begin
    
    parent_comment_id := coalesce(NEW.parent_id, OLD.parent_id);

    
    if parent_comment_id is null then
        return null;
    end if;

   
    select count(*) into reply_count
    from public.comments
    where parent_id = parent_comment_id and deleted_at is null;

    
    update public.comments
    set attributes = jsonb_set(
        attributes,
        '{replies_count}',
        to_jsonb(reply_count),
        true
    )
    where id = parent_comment_id;

    return null; 
end;
$$ language plpgsql;

create trigger update_comment_replies_count_trigger
after insert or delete or update on public.comments
for each row
execute function update_comment_replies_count();