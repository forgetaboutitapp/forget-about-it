CREATE TABLE Logins (
    login_uuid text primary key,
    user_uuid text NOT NULL,
    device_description text NOT NULL,
    FOREIGN KEY (user_uuid) REFERENCES Users(user_uuid)
)

