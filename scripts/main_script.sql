SELECT 'CREATE DATABASE CV'
WHERE NOT EXISTS(SELECT FROM PG_DATABASE WHERE DATNAME = 'CV');
\gexec

CREATE TABLE IF NOT EXISTS Users
(
	Id SERIAL PRIMARY KEY,
	Name VARCHAR UNIQUE NOT NULL,
	Password VARCHAR NOT NULL,
	IsAdmin BOOLEAN DEFAULT FALSE,
	Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	Updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Requests
(
	Id SERIAL PRIMARY KEY,
	UserId INTEGER NOT NULL,
	Status VARCHAR CHECK (
		Status = 'Accepted' OR
		Status = 'Rejected' OR
		Status = 'Submitted'
	) DEFAULT 'Submitted',
	Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	Updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fkUser
		FOREIGN KEY(UserId)
			REFERENCES Users(Id)
);

CREATE TABLE IF NOT EXISTS Candidates
(
	Id SERIAL PRIMARY KEY,
	RequestId INTEGER NOT NULL,
	Name VARCHAR NOT NULL,
	Surname VARCHAR NOT NULL,
	CvOsFileId INTEGER,
	CONSTRAINT fkRequest
		FOREIGN KEY(RequestId)
			REFERENCES Requests(Id)
);
