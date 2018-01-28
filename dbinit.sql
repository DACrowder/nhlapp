create database nhlapp;

\c nhlapp

create table event (
    event_id int PRIMARY KEY,
    event_type text,
    player1_id int,
    player2_id int,
    player1_type text,
    player2_type text,
    player1_team text,
    player2_team text,
    coord_x int,
    coord_y int,
    period int,
    period_time int,
    game_id int,
    UNIQUE (game_id, event_id)
);

create table shift (
    game_id int,
    player_id int,
    period int,
    time_start int,
    time_end int,
    UNIQUE(game_id, player_id, period, time_start, time_end)
);
