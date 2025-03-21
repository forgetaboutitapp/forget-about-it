CREATE TABLE questions(
    question_id integer not null primary key,
    user_id integer not null,
    question text not null,
    answer text not null,
    enabled integer not null default true,
    FOREIGN KEY (user_id) REFERENCES Users(user_id)

);

CREATE TABLE questions_to_tags(
    question_id integer not null,
    tag text not null,
    FOREIGN KEY (question_id) REFERENCES questions(question_id)

);

CREATE TABLE questions_logs(
    question_id integer NOT NULL,
    result integer CHECK( result IN (0, 1) ) NOT NULL,
    timestamp integer NOT NULL,
    FOREIGN KEY (question_id) REFERENCES questions(question_id)
);

CREATE TABLE spacing_algorithms(
    algorithm_id INTEGER PRIMARY KEY,
    alloc text not null,
    api_version integer not null,
    author text not null,
    dealloc text not null,
    desc text,
    download_url text not null,
    init text not null,
    license text not null,
    module_name text not null,
    algorithm_name text not null,
    remote_url text not null,
    version integer not null,
    timestamp integer not null,
    wasm blob not null
);

ALTER TABLE Users ADD default_algorithm INTEGER REFERENCES spacing_algorithms(algorithm_id);