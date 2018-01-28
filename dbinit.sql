drop database if exists nhlapp;

create database nhlapp;

\c nhlapp

create table event (
    event_id int,
    event_type text,
    player1_id int,
    player2_id int,
    player1_type text,
    player2_type text,
    player1_team text,
    coord_x float,
    coord_y float,
    period int,
    period_time int,
    game_id int,
    PRIMARY KEY (game_id, event_id)
);

create table shift (
    game_id int,
    player_id int,
    period int,
    time_start int,
    time_end int,
    team text,
    player_pos text,
    PRIMARY KEY (game_id, player_id, period, time_start, time_end)
);

create table event_roster (
    game_id int,
    event_id int,
    team text,
    player_id int,
    UNIQUE (game_id, event_id, team, player_id)
);
