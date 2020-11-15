CREATE TABLE User(
	isUser INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	userName CHAR(25) NOT NULL,
	EmailUtente CHAR(50) NOT NULL,
	PasswordUtente VARCHAR(15) NOT NULL
);

CREATE TABLE Person(
	idPerson REFERENCES User(idUser),
	name CHAR(20) NOT NULL,
	surname CHAR(20) NOT NULL,
	parents REFERENCES Couple(idCouple)
	
);

CREATE TABLE Couple(
	idCouple INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	father REFERENCES Person(idPerson),
	mother REFERENCES Person(idPerson)
);
