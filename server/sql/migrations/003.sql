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
    question_id integer primary key,
    result integer CHECK( result IN ('correct', 'wrong') ) NOT NULL,
    timestamp integer,
    FOREIGN KEY (question_id) REFERENCES questions(question_id)
);