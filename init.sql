-- Creation of product table

CREATE TABLE IF NOT EXISTS personalDeduction (
	id SERIAL PRIMARY KEY,
	amount INT NOT NULL
);

CREATE TABLE IF NOT EXISTS Kreceipt (
	id SERIAL PRIMARY KEY,
	amount INT NOT NULL
);

INSERT INTO personalDeduction (amount) VALUES (1000);
INSERT INTO Kreceipt (amount) VALUES (1000);