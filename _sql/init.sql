CREATE TABLE users  (
     id text primary key,
     balance int,
);

CREATE TABLE orders (
     id_user text,
     id_service text,
     id_order text,
     amount int,
     accepted boolean,
     primary key (id_user, id_service, id_order),
     foreign key (id_user) references users (id)
);

CREATE INDEX on users (id);
CREATE INDEX on orders (id_user, id_service, id_order);