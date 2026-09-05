CREATE SCHEMA booking;


CREATE TABLE booking.users (
    id           SERIAL       PRIMARY KEY,
    version      BIGINT                     NOT NULL    DEFAULT 1,
    name         VARCHAR(100)               NOT NULL    CHECK(char_length(name) BETWEEN 3 AND 100),
    phone_number VARCHAR(15)  UNIQUE                    CHECK(
                                                                phone_number ~ '^\+[0-9]+$'
                                                                AND
                                                                char_length(phone_number) BETWEEN 10 AND 15),
    mail         VARCHAR(100) UNIQUE                    CHECK(
                                                                mail ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
                                                                AND
                                                                char_length(mail) BETWEEN 9 AND 100),
    created_at   TIMESTAMP                  NOT NULL
);


CREATE TABLE booking.events (
    id              SERIAL          PRIMARY KEY,
    title           VARCHAR(500)                                    NOT NULL,
    description     TEXT,
    total_tickets   INTEGER                                         NOT NULL,
    created_at      TIMESTAMP                                       NOT NULL,
    event_place     VARCHAR(1000)                                   NOT NULL,
    event_time      TIMESTAMP                                       NOT NULL,
    owners_user_id  INTEGER         REFERENCES booking.users(id)    NOT NULL
);


CREATE TABLE booking.bookings (
    id          SERIAL      PRIMARY KEY,
    event_id    INTEGER     REFERENCES booking.events(id)   NOT NULL,
    user_id     INTEGER     REFERENCES booking.users(id)    NOT NULL,
    created_at  TIMESTAMP                                   NOT NULL   DEFAULT now(),
    canceled    BOOLEAN                                     NOT NULL   DEFAULT FALSE
);