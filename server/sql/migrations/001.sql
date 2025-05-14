create table  Users
(
    user_id           integer primary key,
    role              integer not null,
    created           integer not null,
    default_algorithm integer,
    foreign key (default_algorithm) references spacing_algorithms(algorithm_id)  deferrable initially deferred
);

create table Logins
(
    login_uuid         text primary key,
    user_id            integer not null,
    device_description text    not null,
    created            integer not null,
    index_id           integer not null,

    foreign key (user_id) references Users (user_id) deferrable initially deferred
);

create table Logs_Logins
(
    login_uuid   text    not null,
    current_time integer not null,
    foreign key (login_uuid) references Logins (login_uuid) deferrable initially deferred
);

create table questions
(
    question_id integer not null primary key,
    user_id     integer not null,
    question    text    not null,
    answer      text    not null,
    explanation text    not null,
    memo_hint   text    not null,
    enabled     integer not null default true,
    foreign key (user_id) references Users (user_id) deferrable initially deferred

);

create table questions_to_tags
(
    question_id integer not null,
    tag         text    not null,
    foreign key (question_id) references questions (question_id) deferrable initially deferred
);

create table questions_logs
(
    question_id integer                            not null,
    result      integer check ( result IN (0, 1) ) not null,
    timestamp   integer                            not null,
    foreign key (question_id) references questions (question_id) deferrable initially deferred
);

create table  spacing_algorithms
(
    algorithm_id   integer primary key,
    alloc          text    not null,
    api_version    integer not null,
    author         text    not null,
    dealloc        text    not null,
    desc           text,
    download_url   text    not null,
    init           text    not null,
    license        text    not null,
    module_name    text    not null,
    algorithm_name text    not null,
    remote_url     text    not null,
    version        integer not null,
    timestamp      integer not null,
    wasm           blob    not null
);
