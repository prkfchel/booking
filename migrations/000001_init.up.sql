CREATE SCHEMA booking;


CREATE TABLE booking.users (
    id              SERIAL          PRIMARY KEY,
    version         BIGINT                      NOT NULL        DEFAULT 1,
    name            VARCHAR(100)                NOT NULL        CHECK(char_length(name) BETWEEN 3 AND 100),
    phone_number    VARCHAR(15)     UNIQUE                      CHECK(
                                                                    phone_number ~ '^\+[0-9]+$'
                                                                    AND
                                                                    char_length(phone_number) BETWEEN 10 AND 15),
    mail            VARCHAR(100)    UNIQUE                      CHECK(
                                                                    mail ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
                                                                    AND
                                                                    char_length(mail) BETWEEN 9 AND 100),
    created_at      TIMESTAMPTZ                 NOT NULL,
    deleted_at      TIMESTAMPTZ                                 DEFAULT NULL,
    CHECK(
        (deleted_at IS NULL)
        OR
        (deleted_at IS NOT NULL AND deleted_at > created_at)
    )
);


CREATE TABLE booking.events (
    id              SERIAL          PRIMARY KEY,
    version         BIGINT                                          NOT NULL    DEFAULT 1,
    title           VARCHAR(500)                                    NOT NULL    CHECK(char_length(title) BETWEEN 2 AND 500),
    description     TEXT                                                        CHECK(char_length(description) > 5),
    total_tickets   INTEGER                                         NOT NULL    CHECK(total_tickets > 0),
    created_at      TIMESTAMPTZ                                     NOT NULL    DEFAULT now(),
    event_address   VARCHAR(1000)                                   NOT NULL,
    event_time      TIMESTAMP                                       NOT NULL,
    owners_user_id  INTEGER         REFERENCES booking.users(id)    NOT NULL,
    canceled        BOOLEAN                                         NOT NULL    DEFAULT FALSE,
    canceled_at     TIMESTAMPTZ                                                   DEFAULT NULL,
    CHECK( 
        (canceled = FALSE AND canceled_at IS NULL)
        OR 
        (canceled AND canceled_at IS NOT NULL AND canceled_at > created_at)
    )
);


CREATE TABLE booking.bookings (
    id          SERIAL      PRIMARY KEY,
    event_id    INTEGER     REFERENCES booking.events(id)   NOT NULL,
    user_id     INTEGER     REFERENCES booking.users(id)    NOT NULL,
    created_at  TIMESTAMPTZ                                 NOT NULL     DEFAULT now(),
    canceled    BOOLEAN                                     NOT NULL     DEFAULT FALSE,
    canceled_at TIMESTAMPTZ                                              DEFAULT NULL
    CHECK( 
        (canceled = FALSE AND canceled_at IS NULL)
        OR 
        (canceled AND canceled_at IS NOT NULL AND canceled_at >= created_at)
    )
);