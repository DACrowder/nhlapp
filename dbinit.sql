create database nhlapp;

\c nhlapp

create table event (
    event_id int PRIMARY KEY,
    event_type text,
    player1_id int,
    player2_id int,
    player1_type text,
    player2_type text,
    coord_x int,
    coord_y int
);

create table game (
    game_id int PRIMARY KEY,
    event_id int REFERENCES event (event_id),
    event_time timestamp
);

create table shift (
    player_id int,
    time_start timestamp,
    time_end timestamp,
    UNIQUE(player_id, time_start, time_end)
);
