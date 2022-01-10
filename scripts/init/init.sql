-- inserting users
INSERT INTO
	users(name, password)
VALUES
	('user1', 'password1'),
	('user2', 'password2');

-- inserting admin
INSERT INTO
	users(name, password, isadmin)
VALUES
	('admin', 'admin', true);

-- inserting users requests
INSERT INTO
	requests(author_id, candidate_name, candidate_surname, cv_file_id)
VALUES
	(1, 'Name', 'Surname', 'something_similar_to_hash'),
	(1, 'Billie', 'Jean', 'something_similar_to_hash'),
	(2, 'Igor', 'Nikolaev', 'something_similar_to_hash'),
	(2, 'Joseph', 'Stalin', 'something_similar_to_hash'),
	(2, 'Bjarne', 'Stroustrup', 'something_similar_to_hash'),