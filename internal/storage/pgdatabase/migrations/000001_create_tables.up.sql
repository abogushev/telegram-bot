
create table currencies(
    code varchar(10) PRIMARY KEY,
    ratio decimal(6) not null
);

create table state(
    current_currency_code varchar REFERENCES currencies (code)
);

create table categories(
name varchar(100) PRIMARY KEY
);

create table spendings(
    value decimal(100, 2) not null,
    category varchar REFERENCES categories (name),
    date date not null
);

--
insert into currencies(code, ratio) values('rub', 1);
insert into state values('rub');
insert into categories(name) values('food'),('other');