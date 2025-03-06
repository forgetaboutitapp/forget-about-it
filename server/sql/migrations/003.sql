CREATE TABLE questions(
    question_id integer not null primary key,
    user_id integer not null,
    question string not null,
    answer string not null,
    FOREIGN KEY (user_id) REFERENCES Users(user_id)

);

CREATE TABLE questions_to_tags(
    question_id integer,
    tag string,
    FOREIGN KEY (question_id) REFERENCES questions(question_id)

);

CREATE TABLE questions_logs(
    question_id integer primary key,
    result integer CHECK( result IN ('correct', 'wrong') ) NOT NULL,
    timestamp integer,
    FOREIGN KEY (question_id) REFERENCES questions(question_id)
);