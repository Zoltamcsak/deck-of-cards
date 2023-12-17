create table if not exists decks (
    id varchar(50) primary key,
    shuffled bool default false not null,
    remaining int not null,
    cards text[] not null,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp default current_timestamp not null
)