CREATE TABLE IF NOT EXISTS booking
(
    id         VARCHAR(36) PRIMARY KEY,
    stock_id   VARCHAR(36) NOT NULL,
    user_id    VARCHAR(36) NOT NULL,
    order_id   VARCHAR(36) NOT NULL,
    count      INT         NOT NULL,
    created_at TIMESTAMP   NOT NULL,
    expired_at TIMESTAMP   NOT NULL
);

CREATE TABLE IF NOT EXISTS stock
(
    id           VARCHAR(36) PRIMARY KEY,
    event_id     VARCHAR(36) NOT NULL,
    ticket_id    VARCHAR(36) NOT NULL,
    seats_total  INT         NOT NULL,
    seats_booked INT         NOT NULL
);

ALTER TABLE booking ADD FOREIGN KEY (stock_id) REFERENCES stock (id);

CREATE UNIQUE INDEX IF NOT EXISTS event_id_ticket_id ON stock (event_id, ticket_id);