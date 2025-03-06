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

-- name: FindLoginIDByUser :many
SELECT Logins.device_description, Logins.created, Logins.login_uuid, max(Logs_Logins.current_time) as lastUsed FROM Logins LEFT OUTER JOIN Logs_Logins ON Logins.login_uuid=Logs_Logins.login_uuid WHERE Logins.user_uuid=? GROUP BY Logins.device_description, Logins.login_uuid;

-- name: RegisterLogin :exec
INSERT INTO Logs_Logins(login_uuid, current_time) VALUES(?, ?);
