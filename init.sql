CREATE TABLE transactions (
  "user" VARCHAR(20),
  "transaction" NUMERIC(12, 2),
  balance NUMERIC(12, 2) DEFAULT 0,
  date TIMESTAMP
);

INSERT INTO transactions ("user", "transaction", balance, date) VALUES
('bamboo', 0, 1500, '2024-02-25 13:00'),
('kawayan', 0, 2000, '2024-02-25 13:00'),
('hunyo', 0, 2500, '2024-02-25 13:00');
