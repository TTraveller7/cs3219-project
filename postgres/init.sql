DROP TABLE IF EXISTS users, difficulty, matches, questions;

CREATE TABLE users (
    username TEXT PRIMARY KEY,
    password TEXT NOT NULL
);

CREATE TABLE difficulty (
    difficulty TEXT PRIMARY KEY
);

CREATE TABLE matches (
    match_id TEXT PRIMARY KEY,
    username_a TEXT REFERENCES users(username),
    username_b TEXT REFERENCES users(username),
    difficulty TEXT REFERENCES difficulty(difficulty),
    is_ended TEXT,
    created_at TIMESTAMP,
    ended_at TIMESTAMP
);

begin;
    INSERT INTO users 
    VALUES ('Bob', '$2a$10$WluXDB/sWvvW5AmseUwXKebuaOCHvS3gEERviX9jekX2A9ibQi8nm'), ('Felix', '$2a$10$2B31v9TDfs3wH85tTHJWUeYqLaRe4V8widBuqTFAcD0a88U1epFyi'); 

    INSERT INTO difficulty
    VALUES ('easy'), ('medium'), ('hard');
commit;
