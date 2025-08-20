-- для самих каналов
create table feeds (
    id serial primary key,
    title text not null,
    link text not null,
    description text
);

-- для самих статей
create table items (
    id serial primary key,
    feed_id int references feeds(id) on delete cascade,
    title text not null,
    link text not null,
    description text,
    pub_date timestamp
);