-- inserting users
INSERT INTO USERS(NAME, PASSWORD) VALUES ('user1', 'password1'), ('user2', 'password2');

-- inserting admin
INSERT INTO USERS(NAME, PASSWORD, ISADMIN) VALUES ('admin', 'admin', TRUE);

-- inserting candidates
INSERT INTO CANDIDATES(NAME, SURNAME, CVOSFILEID)
VALUES ('can1', 'didate1', '1111'),
		('can2', 'didate2', '1111'),
		('can3', 'didate3', '1111'),
		('can4', 'didate4', '1111'),
		('can5', 'didate5', '1111'),
		('can6', 'didate6', '1111'),
		('can7', 'didate7', '1111');

-- inserting users requests
INSERT INTO REQUESTS(USERID, CANDIDATEID, STATUS) VALUES (1,2, 'Rejected'), (1,5 ,'Submitted'), (1,6,'Submitted'), (2,1,'Submitted'), (2,3, 'Rejected'), (2,4, 'Rejected'), (1,7, 'Submitted');
