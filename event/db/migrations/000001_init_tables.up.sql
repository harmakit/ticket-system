CREATE TABLE IF NOT EXISTS event
(
    id          VARCHAR(36) PRIMARY KEY,
    start_date  TIMESTAMP    NOT NULL,
    duration    SMALLINT     NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NULL,
    location_id VARCHAR(36)  NULL
);

CREATE TABLE IF NOT EXISTS ticket
(
    id       UUID PRIMARY KEY,
    event_id VARCHAR(36)  NOT NULL,
    name     VARCHAR(255) NOT NULL,
    price    DECIMAL      NOT NULL,
    FOREIGN KEY (event_id) REFERENCES event (id)
);

CREATE TABLE IF NOT EXISTS location
(
    id      VARCHAR(36) PRIMARY KEY,
    name    VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    lat     DECIMAL      NOT NULL,
    lng     DECIMAL      NOT NULL
);

ALTER TABLE event ADD FOREIGN KEY (location_id) REFERENCES location (id);