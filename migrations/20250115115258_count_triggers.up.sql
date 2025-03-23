create or replace function update_user_education_count()
returns trigger as $$
declare
    count integer;
begin
    select count(*) into count from public.educations where user_id = NEW.user_id;

    update public.user_profiles
    set attributes = jsonb_set(attributes, '{education_metadata}', coalesce(attributes -> 'education_metadata', '{}'::jsonb) || jsonb_build_object('count', count), true)
    where user_id = NEW.user_id;   
    
    return NEW;
end;
$$ language plpgsql;

create trigger update_user_education_count_trigger
after insert or delete on public.educations
for each row
execute function update_user_education_count();



create or replace function update_user_experience_count()
returns trigger as $$
declare
    count integer;
begin
    select count(*) into count from public.work_experiences where user_id = NEW.user_id;

    update public.user_profiles
    set attributes = jsonb_set(attributes, '{work_experience_metadata}', coalesce(attributes -> 'work_experience_metadata', '{}'::jsonb) || jsonb_build_object('count', count), true)
    where user_id = NEW.user_id;   
    
    return NEW;
end;
$$ language plpgsql;

create trigger update_user_experience_count_trigger
after insert or delete on public.work_experiences
for each row
execute function update_user_experience_count();



create or replace function update_user_certification_count()
returns trigger as $$
declare
    count integer;
begin
    select count(*) into count from public.certifications where user_id = NEW.user_id;

    update public.user_profiles
    set attributes = jsonb_set(attributes, '{certification_metadata}', coalesce(attributes -> 'certification_metadata', '{}'::jsonb) || jsonb_build_object('count', count), true)
    where user_id = NEW.user_id;   
    
    return NEW;
end;
$$ language plpgsql;

create trigger update_user_certification_count_trigger
after insert or delete on public.certifications
for each row
execute function update_user_certification_count();



create or replace function update_user_hackathon_count()
returns trigger as $$
declare
    count integer;
begin
    select count(*) into count from public.hackathons where user_id = NEW.user_id;

    update public.user_profiles
    set attributes = jsonb_set(attributes, '{hackathon_metadata}', coalesce(attributes -> 'hackathon_metadata', '{}'::jsonb) || jsonb_build_object('count', count), true)
    where user_id = NEW.user_id;   
    
    return NEW;
end;
$$ language plpgsql;

create trigger update_user_hackathon_count_trigger
after insert or delete on public.hackathons
for each row
execute function update_user_hackathon_count();

create or replace function update_user_blog_count()
returns trigger as $$
declare
    count integer;
begin
    select count(blogs.id) into count from public.blogs where user_id = NEW.user_id;

    update public.user_profiles
    set attributes = jsonb_set(attributes, '{blog_metadata}', coalesce(attributes -> 'blog_metadata', '{}'::jsonb) || jsonb_build_object('count', count), true)
    where user_id = NEW.user_id;   
    
    return NEW;
end;
$$ language plpgsql;

create trigger update_user_blogs_count_trigger
after insert or delete on public.blogs
for each row
execute function update_user_blog_count();

create or replace function update_user_work_gallery_count()
returns trigger as $$
declare
    count integer;
begin
    select count(tech_projects.id) into count from public.tech_projects where user_id = NEW.user_id;

    update public.tech_projects
    set attributes = jsonb_set(attributes, '{work_gallery_metadata}', coalesce(attributes -> 'work_gallery_metadata', '{}'::jsonb) || jsonb_build_object('count', count), true)
    where user_id = NEW.user_id;   
    
    return NEW;
end;
$$ language plpgsql;

create trigger update_user_work_gallery_count_trigger
after insert or delete on public.tech_projects
for each row
execute function update_user_work_gallery_count();
