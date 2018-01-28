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
    PRIMARY KEY (game_id, event_id, team, player_id)
);

create table line (
    line_id serial,
    game_id int,
    line_players text,
    event_id int,
    team text,
    PRIMARY KEY (game_id, line_players)
);

create table event_winners (
    line_players text,
    team text,
    event_id int,
    game_id int,
    PRIMARY KEY (line_players, team, event_id, game_id)
);

create table event_losers (
    line_players text,
    team text,
    event_id int,
    game_id int,
    PRIMARY KEY (line_players, team, event_id, game_id)
);

create table lineups (
    line_players text,
    team text,
    game_id int,
    hits int,
    missed_shot int,
    blocked_shot int,
    shot int,
    goal int,
    faceoff int,
    takeaway int,
    penalty int,
    giveaway int,
    PRIMARY KEY (line_players, game_id)
);
