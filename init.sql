-- Creation of product table

CREATE TABLE IF NOT EXISTS ktaxes (
	id SERIAL PRIMARY KEY,
	amount FLOAT NOT NULL,
	taxType VARCHAR(255) NOT NULL
);


INSERT INTO ktaxes (amount, taxType) VALUES (50000.0, 'k-receipt');
INSERT INTO ktaxes (amount, taxType) VALUES (60000.0, 'personalDeduction');