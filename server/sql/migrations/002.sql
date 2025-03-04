CREATE TABLE Logins (
    login_uuid text primary key,
    user_uuid text NOT NULL,
    device_description text NOT NULL,
    created int NOT NULL,
    FOREIGN KEY (user_uuid) REFERENCES Users(user_uuid)
);

CREATE TABLE Logs_Logins (
    login_uuid text,
    current_time int,
    FOREIGN KEY (login_uuid) REFERENCES Logins(login_uuid)
);