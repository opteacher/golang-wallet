# CreateTable
CREATE TABLE IF NOT EXISTS collect (
  id INTEGER NOT NULL AUTO_INCREMENT,
  tx_hash VARCHAR(255) UNIQUE,
  address VARCHAR(255) NOT NULL,
  amount DECIMAL(64,20) NOT NULL,
  asset CHAR(32) NOT NULL,
  create_time DATETIME DEFAULT NOW(),
  PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

# AddSentCollect
INSERT INTO collect (%s address, amount, asset) VALUES (%s ?, ?, ?)