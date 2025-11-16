-- +migrate Up

CREATE TABLE teams (
    team_name VARCHAR(64) PRIMARY KEY
);

CREATE TABLE users (
    user_id VARCHAR(64) PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    team_name VARCHAR(64) NOT NULL REFERENCES teams(team_name)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL
);

CREATE TYPE status_enum AS ENUM ('OPEN', 'MERGED');

CREATE TABLE pull_requests (
    pull_request_id VARCHAR(64) PRIMARY KEY,
    pull_request_name VARCHAR(64) NOT NULL,
    author_id VARCHAR(64) NOT NULL REFERENCES users(user_id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    status status_enum NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE pr_reviewers (
    pr_reviewers_id SERIAL PRIMARY KEY,
    pull_request_id VARCHAR(64) NOT NULL REFERENCES pull_requests(pull_request_id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    user_id VARCHAR(64) NOT NULL REFERENCES users(user_id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT pr_reviewers_pull_request_id_user_id_unique UNIQUE (pull_request_id, user_id)
);
