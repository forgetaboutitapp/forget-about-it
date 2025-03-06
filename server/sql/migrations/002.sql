CREATE TABLE Logins (
    login_uuid text primary key,
    user_id integer NOT NULL,
    device_description text NOT NULL,
    created integer NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

CREATE TABLE Logs_Logins (
    login_uuid text not null,
    current_time integer not null,
    FOREIGN KEY (login_uuid) REFERENCES Logins(login_uuid)
);