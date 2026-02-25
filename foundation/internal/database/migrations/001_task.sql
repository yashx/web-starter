-- +goose Up
CREATE TABLE task
(
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    version     INT          NOT NULL DEFAULT 0
);

INSERT INTO task (id, description, created_at, updated_at, version)
VALUES ('c716d20a-6214-4125-a30c-78dac2f3139d', 'Walk', NOW(), NOW(), 0);