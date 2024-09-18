CREATE TABLE IF NOT EXISTS "order"
(
    id      VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36)  NOT NULL,
    status  VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS item
(
    id       VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    stock_id VARCHAR(36) NOT NULL,
    count    INT         NOT NULL
);

CREATE TABLE IF NOT EXISTS cart
(
    id        VARCHAR(36) PRIMARY KEY,
    user_id   VARCHAR(36) NOT NULL,
    ticket_id VARCHAR(36) NOT NULL,
    count     INT         NOT NULL
);

ALTER TABLE item ADD FOREIGN KEY (order_id) REFERENCES "order" (id);

CREATE UNIQUE INDEX IF NOT EXISTS user_id_ticket_id ON cart (user_id, ticket_id);