create sequence appusers_id_seq;

alter sequence appusers_id_seq owner to ramazan;

create sequence notes_id_seq;

alter sequence notes_id_seq owner to ramazan;

create table ramazan
(
    name text
);

alter table ramazan owner to ramazan;

create table appusers
(
    id serial not null
        constraint appusers_pk
            primary key,
    login varchar,
    password integer not null,
    token varchar
);

alter table appusers owner to ramazan;

create unique index appusers_login_uindex
    on appusers (login);

create table notes
(
    id serial not null
        constraint notes_pk
            primary key,
    name varchar,
    text varchar,
    date timestamp,
    done boolean,
    color varchar,
    token varchar
);

alter table notes owner to ramazan;