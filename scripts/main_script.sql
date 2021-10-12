CREATE TABLE Users
(
	ID SERIAL PRIMARY KEY,
	Name VARCHAR NOT NULL,
	Password VARCHAR NOT NULL,
	Type VARCHAR CHECK (
		Type = 'Admin' OR 
		Type = 'User'
	)
);

CREATE TABLE Candidates
(
	ID SERIAL PRIMARY KEY,
	User_ID INTEGER NOT NULL,
	Name VARCHAR NOT NULL,
	Surname VARCHAR NOT NULL,
	Publish_date DATE,
	CONSTRAINT fk_user
		FOREIGN KEY(User_ID)
			REFERENCES Users(ID)
);

CREATE TABLE CVs
(
	ID SERIAL PRIMARY KEY,
	Candidate_ID INTEGER NOT NULL,
	Status VARCHAR CHECK (
		Status = 'Accepted' OR 
		Status = 'Rejected' OR 
		Status = 'Submitted' OR 
		Status = 'Pending'
	),
	OS_File_ID INTEGER,
	CONSTRAINT fk_candidate
		FOREIGN KEY(Candidate_ID)
			REFERENCES Candidates(ID)
);