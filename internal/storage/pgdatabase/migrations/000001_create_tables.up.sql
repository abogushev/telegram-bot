
create table currencies(
    code varchar(10) PRIMARY KEY,
    ratio decimal(10, 6) not null
);

create table state(
    current_currency_code varchar REFERENCES currencies (code)
);

create table categories(
    id integer PRIMARY KEY,
    name varchar(100) unique
);

create table spendings(
    value decimal(100, 2) not null,
    category_id INTEGER REFERENCES categories (id),
    date date not null
);

--
insert into currencies(code, ratio) values('rub', 1);
insert into state values('rub');
insert into categories(id, name) values(0, 'food'),(1, 'other');