create or replace function update_user_education_count()
returns trigger as $$
declare
    count integer;
begin
    select count(*) into count from public.user_educations where user_id = NEW.user_id;

    update public.user_profiles
    set attributes = jsonb_set(attributes, '{education_metadata}', coalesce(attributes -> 'education_metadata', '{}'::jsonb) || jsonb_build_object('count', count), true)
    where user_id = NEW.user_id;   
    
    return NEW;
end;
$$ language plpgsql;

create trigger update_user_education_count_trigger
after insert or delete on public.user_educations
for each row
execute function update_user_education_count();



create or replace function update_user_experience_count()
returns trigger as $$
declare
    count integer;
begin
    select count(*) into count from public.user_experiences where user_id = NEW.user_id;

    update public.user_profiles
    set attributes = jsonb_set(attributes, '{work_experience_metadata}', coalesce(attributes -> 'work_experience_metadata', '{}'::jsonb) || jsonb_build_object('count', count), true)
    where user_id = NEW.user_id;   
    
    return NEW;
end;
$$ language plpgsql;

create trigger update_user_experience_count_trigger
after insert or delete on public.user_experiences
for each row
execute function update_user_experience_count();



create or replace function update_user_certification_count()
returns trigger as $$
declare
    count integer;
begin
    select count(*) into count from public.user_certifications where user_id = NEW.user_id;

    update public.user_profiles
    set attributes = jsonb_set(attributes, '{certification_metadata}', coalesce(attributes -> 'certification_metadata', '{}'::jsonb) || jsonb_build_object('count', count), true)
    where user_id = NEW.user_id;   
    
    return NEW;
end;
$$ language plpgsql;

create trigger update_user_certification_count_trigger
after insert or delete on public.user_certifications
for each row
execute function update_user_certification_count();



create or replace function update_user_hackathon_count()
returns trigger as $$
declare
    count integer;
begin
    select count(*) into count from public.user_hackathons where user_id = NEW.user_id;

    update public.user_profiles
    set attributes = jsonb_set(attributes, '{hackathon_metadata}', coalesce(attributes -> 'hackathon_metadata', '{}'::jsonb) || jsonb_build_object('count', count), true)
    where user_id = NEW.user_id;   
    
    return NEW;
end;
$$ language plpgsql;

create trigger update_user_hackathon_count_trigger
after insert or delete on public.user_hackathons
for each row
execute function update_user_hackathon_count();
