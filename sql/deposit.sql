# CreateTable
CREATE TABLE IF NOT EXISTS deposit (
  id INTEGER NOT NULL AUTO_INCREMENT,
  tx_hash VARCHAR(255) NOT NULL,
  address VARCHAR(255) NOT NULL,
  amount DECIMAL(30,20) NOT NULL,
  asset CHAR(32) NOT NULL,
  height INTEGER(11) NOT NULL,
  tx_index INTEGER,
  create_time DATETIME DEFAULT NOW(),
  update_time DATETIME DEFAULT NOW() ON UPDATE NOW(),
  PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

# FirstFindDeposit
INSERT INTO deposit (tx_hash, address, amount, asset, height, tx_index) VALUES (?, ?, ?, ?, ?, ?)