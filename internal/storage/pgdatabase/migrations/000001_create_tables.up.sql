
create table currencies(
    code varchar(10) PRIMARY KEY,
    ratio decimal(10, 6) not null
);

create table state(
    current_currency_code varchar REFERENCES currencies (code),
    budget_value decimal(10, 2) not null,
    budget_balance decimal(10, 2) not null,
    budget_expires_in date not null
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

-- индекс по дате, т.к по ней происходит отбор для отчета
CREATE INDEX idx_spendings_date ON spendings(date);
--
insert into currencies(code, ratio) values('rub', 1);
insert into categories(id, name) values(0, 'food'),(1, 'other');

insert into state(current_currency_code, budget_value, budget_balance, budget_expires_in) values('rub', 1000, 1000, now() + interval '1 month');