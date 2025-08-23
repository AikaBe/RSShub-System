-- для самих каналов
create table feeds (
    id serial primary key,
    name text not null unique,
    url text not null unique,
    created_at timestamp not null ,
    updated_at timestamp not null
);

-- для самих статей
create table articles (
    id serial primary key,
    feed_id int references feeds(id) on delete cascade,
    title text not null,
    link text not null,
    published_at timestamp,
    description text not null ,
    created_at timestamp,
    updated_at timestamp
);