DROP TABLE IF EXISTS report;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;


CREATE TABLE users  (
     id text primary key,
     balance int
);

CREATE TABLE orders (
     id_user text,
     id_service text,
     id_order text UNIQUE,
     amount int,
     accepted boolean,
     primary key (id_user, id_service, id_order),
     foreign key (id_user) references users (id)
);

CREATE TABLE report (
     id_user text,
     id_service text,
     id_order text,
     amount int,
     accepted_at date,
     primary key (id_user, id_service, id_order),
     foreign key (id_order) references orders (id_order)
);

CREATE INDEX on users (id);
CREATE INDEX on orders (id_user, id_service, id_order);
CREATE INDEX on report (accepted_at);