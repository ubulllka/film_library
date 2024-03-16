CREATE TABLE IF NOT EXISTS users(
    id                  SERIAL PRIMARY KEY  NOT NULL UNIQUE,
    username            VARCHAR(255)        NOT NULL UNIQUE,
    user_role           VARCHAR(255)        NOT NULL,
    password_hash       VARCHAR(255)        NOT NULL
);

CREATE TABLE IF NOT EXISTS actors(
    id                  SERIAL PRIMARY KEY  NOT NULL UNIQUE,
    actor_name          VARCHAR(255)        NOT NULL,
    sex                 VARCHAR(255)        NOT NULL,
    b_date              DATE                NOT NULL
);

CREATE TABLE IF NOT EXISTS films(
    id                  SERIAL PRIMARY KEY  NOT NULL UNIQUE,
    film_name           VARCHAR(255)        NOT NULL,
    description         TEXT                NOT NULL,
    release_date        DATE                NOT NULL,
    rating              FLOAT               NOT NULL
);

CREATE TABLE IF NOT EXISTS film_actor(
    id          SERIAL PRIMARY KEY                          NOT NULL UNIQUE,
    film_id     INT REFERENCES films(id) ON DELETE CASCADE  NOT NULL,
    actor_id    INT REFERENCES actors(id) ON DELETE CASCADE NOT NULL
);


