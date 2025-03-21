-- name: AddUser :exec
INSERT INTO Users(user_id, role, created) VALUES(?, ?, ?);

-- name: AddLogin :exec
INSERT INTO Logins(login_uuid, user_id, device_description, created) VALUES(?, ?, ?, ?);

-- name: GetUser :many
SELECT user_id FROM Users where role=?;

-- name: FindAdminLogins :many
SELECT DISTINCT Logins.user_id FROM Logins JOIN Users ON Logins.user_id = Users.user_id where Users.role=0;

-- name: FindUserByLogin :many
SELECT DISTINCT Users.user_id FROM Users JOIN Logins ON Logins.user_id = Users.user_id WHERE Logins.login_uuid = ?;

-- name: SetLastLogin :exec
INSERT INTO Logs_Logins (login_uuid, current_time) VALUES (?, ?);

-- name: CreateNewLogin :one
INSERT INTO Logins (login_uuid, user_id, device_description, created) values (?, ?, ?, ?) returning login_uuid;

-- name: FindLoginIDByUser :many
SELECT Logins.device_description, Logins.created, Logins.login_uuid, max(Logs_Logins.current_time) as lastUsed FROM Logins LEFT OUTER JOIN Logs_Logins ON Logins.login_uuid=Logs_Logins.login_uuid WHERE Logins.user_id=? GROUP BY Logins.device_description, Logins.login_uuid;

-- name: RegisterLogin :exec
INSERT INTO Logs_Logins(login_uuid, current_time) VALUES(?, ?);

-- name: GetAllQuestions :many
SELECT questions.question_id, questions.question, questions.answer FROM questions WHERE questions.user_id=? AND questions.enabled = 1;

-- name: GetTagsByQuestion :many
SELECT tag FROM questions_to_tags WHERE question_id = ?;

-- name: AddNewQuestion :exec
INSERT INTO questions(question_id, user_id, question, answer, enabled) VALUES (?, ?, ?, ?, ?);

-- name: AddNewTag :exec
INSERT INTO questions_to_tags(question_id, tag) VALUES(?, ?);

-- name: DeleteAllTags :exec
DELETE FROM questions_to_tags WHERE questions_to_tags.question_id in (SELECT question_id FROM questions WHERE user_id = ?);

-- name: UpdateQuestion :exec
UPDATE questions SET question=?, answer=?, enabled=? WHERE question_id=?;

-- name: GetTagsByUser :many
with tmpValue (id) as (SELECT question_id FROM QUESTIONS WHERE questions.user_id=?) select distinct tag from questions_to_tags, tmpValue where questions_to_tags.question_id=tmpValue.id;

-- name: GradeQuestion :exec
INSERT INTO questions_logs (question_id, result, timestamp) VALUES (?, ?, ?);

-- name: GetAllGrades :many
SELECT questions_logs.question_id, result, timestamp FROM questions_logs JOIN questions ON questions.question_id = questions_logs.question_id WHERE questions.user_id=?;

-- name: AddSpacingAlgorithm :exec
INSERT INTO spacing_algorithms( algorithm_id,
    alloc,
    api_version,
    author,
    dealloc,
    desc,
    download_url,
    init,
    license,
    module_name,
    algorithm_name,
    remote_url,
    version,
    timestamp,
    wasm) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);

-- name: GetSpacingAlgorithms :many
SELECT algorithm_id, 
    alloc,
    api_version,
    author,
    dealloc,
    desc,
    download_url,
    init,
    license,
    module_name,
    algorithm_name,
    remote_url,
    version,
    timestamp,
    wasm FROM spacing_algorithms;

-- name: GetDefaultAlgorithm :one
SELECT default_algorithm from Users WHERE user_id = ?;