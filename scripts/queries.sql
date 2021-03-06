-- all requests
SELECT ID, USERID, CANDIDATEID, STATUS, CREATED, UPDATED
FROM REQUESTS
ORDER BY ID;

-- all candidates
SELECT ID, NAME, SURNAME, CVOSFILEID FROM CANDIDATES;

-- one user requests
SELECT ID, USERID, CANDIDATEID, STATUS, CREATED, UPDATED FROM REQUESTS WHERE USERID = 1;

-- id of file in object storage
SELECT CVOSFILEID FROM CANDIDATES WHERE ID = 1;

-- user password
SELECT PASSWORD FROM USERS WHERE ID = 2;

-- user login
SELECT NAME FROM USERS WHERE NAME = 'user1';

-- user with the most number of requests
SELECT USERID FROM REQUESTS
GROUP BY USERID
ORDER BY COUNT(ID)
DESC LIMIT 1

-- list of users sorted by the number of requests
SELECT USERID FROM REQUESTS
GROUP BY USERID
ORDER BY COUNT(ID)
DESC

-- list of users which have declined requests
SELECT DISTINCT ID FROM USERS
JOIN REQUESTS ON USERS.ID = REQUESTS.USERID
WHERE REQUESTS.STATUS = 'Rejected';

-- list of candidates which have declined requests
SELECT CANDIDATES.ID FROM CANDIDATES
JOIN REQUESTS ON CANDIDATES.ID = REQUESTS.CANDIDATEID
WHERE REQUESTS.STATUS = 'Rejected';

-- user creating 
INSERT INTO USERS(NAME, PASSWORD) VALUES('Name', 'Password');

-- candidate creating 
INSERT INTO CANDIDATES(NAME, SURNAME, CVOSFILEID) VALUES('Name', 'Surname', '1234-5678');

-- request creating
INSERT INTO REQUESTS(USERID, CANDIDATEID, STATUS) VALUES(1,2,'Submitted')

-- request status updating
UPDATE REQUESTS SET STATUS = 'Accepted' WHERE ID = 8

-- full list of user candidates
SELECT CANDIDATES.NAME, CANDIDATES.SURNAME, REQUESTS.STATUS, REQUESTS.CREATED
FROM CANDIDATES JOIN REQUESTS
ON CANDIDATES.ID = REQUESTS.CANDIDATEID JOIN USERS 
ON REQUESTS.USERID = USERS.ID WHERE USERID = 1