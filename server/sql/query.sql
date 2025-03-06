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
