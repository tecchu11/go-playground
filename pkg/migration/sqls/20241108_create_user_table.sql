-- +goose Up
CREATE TABLE users (
    id BINARY(16) NOT NULL PRIMARY KEY COMMENT 'id is user id',
    sub VARCHAR(64) NOT NULL COMMENT 'sub is jwt subject',
    given_name VARCHAR(64) NOT NULL COMMENT 'given_name of user',
    family_name VARCHAR(64) NOT NULL COMMENT 'family_name of user',
    email VARCHAR(255) NOT NULL COMMENT 'email of user',
    email_verified TINYINT(1) NOT NULL COMMENT 'email_verified is whether email is verified',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_sub (sub) COMMENT 'index for jwt subject'
) COMMENT = 'users is user information';


-- +goose Down
DROP TABLE IF EXISTS users;
