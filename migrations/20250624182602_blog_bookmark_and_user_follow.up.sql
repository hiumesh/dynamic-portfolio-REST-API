create table public.blog_bookmarks (
    id bigserial,
    user_id uuid not null,
    blog_id bigint not null,
    constraint blog_bookmarks_pkey primary key (id),
    constraint blog_bookmarks_blog_id_fkey foreign key (blog_id) references public.blogs (id) on delete cascade,
    constraint blog_bookmarks_user_id_fkey foreign key (user_id) references auth.users (id) on delete cascade,
    constraint blog_bookmarks_user_id_and_blog_id_composite_key unique (user_id, blog_id) deferrable initially deferred
) tablespace pg_default;

create table public.user_follows (
    id bigserial,
    follower_id uuid not null,
    following_id uuid not null,
    constraint user_follows_pkey primary key (id),
    constraint user_follows_follower_id_fkey foreign key (follower_id) references auth.users (id) on delete cascade,
    constraint user_follows_following_id_fkey foreign key (following_id) references auth.users (id) on delete cascade,
    constraint user_follows_follower_id_and_following_id_composite_key unique (follower_id, following_id) deferrable initially deferred
) tablespace pg_default;