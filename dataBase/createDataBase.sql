CREATE TABLE usertype
(
    TypeKey INT PRIMARY KEY,
    TypeVal CHAR(7)
);
INSERT INTO usertype
VALUES (1, 'Admin'),
       (2, 'Student'),
       (3, 'Teacher');

CREATE TABLE member
(
    Nickname     CHAR(20),
    Username     CHAR(20) UNIQUE,
    UserPassword CHAR(20),
    UserType     INT,
    FOREIGN KEY (UserType) REFERENCES usertype (TypeKey)
) DEFAULT CHARSET UTF8;

insert into member (Nickname, Username, UserType, UserPassword)
values ("root", "JudgeAdmin", 1, "JudgePassword2022");
insert into member (Nickname, Username, UserType, UserPassword)
values ("test1", "user1", 1, "JudgePassword2022");
insert into member (Nickname, Username, UserType, UserPassword)
values ("test1", "user2", 2, "JudgePassword2022");
insert into member (Nickname, Username, UserType, UserPassword)
values ("test2", "user3", 3, "JudgePassword2022");