-- name: AddUser :exec
INSERT INTO Users(user_uuid, role) VALUES(?, ?);

-- name: AddLogin :exec
INSERT INTO Logins(login_uuid, user_uuid, device_description) VALUES(?, ?, ?);

-- name: GetUser :many
SELECT user_uuid FROM Users where role=?;

-- name: FindAdminLogins :many
SELECT DISTINCT Logins.user_uuid FROM Logins JOIN Users ON Logins.user_uuid = Users.user_uuid where Users.role=0;