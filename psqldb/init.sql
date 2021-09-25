CREATE DATABASE jiradb;
GRANT ALL PRIVILEGES ON DATABASE jiradb TO jirauser;

CREATE TABLE Users(
    UserID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    Name VARCHAR(50) NOT NULL
);

CREATE TABLE Items(
    ItemID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    Name VARCHAR(50) NOT NULL,
    Description VARCHAR(50),
    Size VARCHAR(50),
    URL VARCHAR(50),
    Type VARCHAR(50)
);

CREATE TABLE Logs(
    LogId INT GENERATED ALWAYS AS IDENTITY,
    UserID INT REFERENCES Users,
    ItemID INT REFERENCES Items,
    LogDate DATE NOT NULL,
    Hours FLOAT(2) NOT NULL,
    Notes VARCHAR(50)
);


-- Dummy data

INSERT INTO Users (Name) VALUES ('SomeUser');

INSERT INTO Items (Name, Description, Size, URL, Type) VALUES ('CR-1234', 'Some CR', 'Large', 'https://jira.om.net/browse/CR-1234', 'CR');
INSERT INTO Items (Name, Description, Size, URL, Type) VALUES ('CR-4321', 'Some other CR', 'Small', 'https://jira.om.net/browse/CR-4321', 'CR');
INSERT INTO Items (Name, Description, Size, URL, Type) VALUES ('Meetings/Training', NULL, NULL, NULL, 'Admin');
INSERT INTO Items (Name, Description, Size, URL, Type) VALUES ('Personal Leave', NULL, NULL, NULL, 'Admin');
INSERT INTO Items (Name, Description, Size, URL, Type) VALUES ('Annual Leave', NULL, NULL, NULL, 'Admin');

INSERT INTO Logs (UserID, ItemID, LogDate, Hours) VALUES (1, 1, CURRENT_DATE, 2);
INSERT INTO Logs (UserID, ItemID, LogDate, Hours, Notes) VALUES (1, 3, CURRENT_DATE, 3.5, 'Team Meeting');
INSERT INTO Logs (UserID, ItemID, LogDate, Hours) VALUES (1, 4, CURRENT_DATE, 2.5);