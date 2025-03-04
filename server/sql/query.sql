-- name: AddUser :exec
INSERT INTO Users(user_uuid, role, created) VALUES(?, ?, ?);

-- name: AddLogin :exec
INSERT INTO Logins(login_uuid, user_uuid, device_description, created) VALUES(?, ?, ?, ?);

-- name: GetUser :many
SELECT user_uuid FROM Users where role=?;

-- name: FindAdminLogins :many
SELECT DISTINCT Logins.user_uuid FROM Logins JOIN Users ON Logins.user_uuid = Users.user_uuid where Users.role=0;

-- name: FindUserByLogin :many
SELECT DISTINCT Users.user_uuid FROM Users JOIN Logins ON Logins.user_uuid = Users.user_uuid WHERE Logins.login_uuid = ?;

-- name: SetLastLogin :exec
INSERT INTO Logs_Logins (login_uuid, current_time) VALUES (?, ?);

-- name: CreateNewLogin :one
INSERT INTO Logins (login_uuid, user_uuid, device_description, created) values (?, ?, ?, ?) returning login_uuid;