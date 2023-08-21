create table
  public.waitlist (
    id uuid not null default gen_random_uuid (),
    created_at timestamp with time zone not null default now(),
    mail text not null default ''::text,
    name text not null default ''::text,
    constraint waitlist_pkey primary key (id),
    constraint waitlist_mail_key unique (mail)
  ) tablespace pg_default;