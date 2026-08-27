CREATE SCHEMA booking;


CREATE TABLE booking.users (
    id           SERIAL       PRIMARY KEY,
    name         VARCHAR(100)               NOT NULL,
    phone_number VARCHAR(15)  UNIQUE        CHECK(phone_number ~ '^\+[1-9][0-9]{7,14}$') ,
    mail         VARCHAR(100) UNIQUE        CHECK(mail ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'),
    created_at   TIMESTAMP                  NOT NULL
);


CREATE TABLE booking.events (
    id              SERIAL          PRIMARY KEY,
    title           VARCHAR(500)                                    NOT NULL,
    description     TEXT,
    total_tickets   INTEGER                                         NOT NULL,
    created_at      TIMESTAMP                                       NOT NULL,
    event_place     VARCHAR(1000)                                   NOT NULL,
    event_time      TIMESTAMP                                       NOT NULL   DEFAULT now(),
    owners_user_id  INTEGER         REFERENCES booking.users(id)    NOT NULL
);


CREATE TABLE booking.bookings (
    id          SERIAL      PRIMARY KEY,
    event_id    INTEGER     REFERENCES booking.events(id)   NOT NULL,
    user_id     INTEGER     REFERENCES booking.users(id)    NOT NULL,
    created_at  TIMESTAMP                                   NOT NULL   DEFAULT now(),
    canceled    BOOLEAN                                     NOT NULL   DEFAULT FALSE
);