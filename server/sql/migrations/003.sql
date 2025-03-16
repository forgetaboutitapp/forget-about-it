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
    author_name TEXT NOT NULL,
    author TEXT NOT NULL,
    license TEXT NOT NULL,
    remote_url TEXT NOT NULL,
    download_url TEXT NOT NULL,
    timestamp_added INTEGER NOT NULL,
    initialization_functions TEXT NOT NULL,
    allocating_function TEXT NOT NULL,
    freeing_function TEXT NOT NULL,
    algorithm BLOB NOT NULL,
    module_name TEXT NOT NULL,
    version INT NOT NULL
    );

ALTER TABLE Users ADD default_algorithm INTEGER REFERENCES spacing_algorithms(algorithm_id);