BEGIN;

-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE test_table (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

COMMIT;