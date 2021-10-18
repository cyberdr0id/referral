-- all users
SELECT * FROM USERS;

-- all requests
SELECT * FROM REQUESTS;

-- all candidates
SELECT * FROM CANDIDATES;

-- one user requests
SELECT * FROM REQUESTS WHERE USERID = 1;

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
