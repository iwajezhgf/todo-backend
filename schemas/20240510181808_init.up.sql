CREATE TABLE users
(
    id       BIGINT AUTO_INCREMENT,
    email    VARCHAR(255) NOT NULL UNIQUE,
    password VARBINARY(255) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE tokens
(
    id      BIGINT AUTO_INCREMENT,
    token   VARCHAR(255) UNIQUE,
    expire  DATETIME,
    user_id BIGINT,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE todos
(
    id      BIGINT AUTO_INCREMENT,
    title   VARCHAR(255) NOT NULL,
    note    TEXT         NOT NULL,
    created DATETIME     NOT NULL,
    expire  DATETIME     NOT NULL,
    status  VARCHAR(15)  NOT NULL DEFAULT "active",
    user_id BIGINT,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
