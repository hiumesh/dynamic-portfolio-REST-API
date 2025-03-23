drop trigger if exists update_user_education_count_trigger on public.educations;
drop function if exists update_user_education_count;

drop trigger if exists update_user_certification_count_trigger on public.certifications;
drop function if exists update_user_certification_count;

drop trigger if exists update_user_hackathon_count_trigger on public.hackathons;
drop function if exists update_user_hackathon_count;

drop trigger if exists update_user_experience_count_trigger on public.work_experiences;
drop function if exists update_user_experience_count;