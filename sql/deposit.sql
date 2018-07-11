# CreateTable
CREATE TABLE IF NOT EXISTS deposit (
  id INTEGER NOT NULL AUTO_INCREMENT,
  tx_hash VARCHAR(255) NOT NULL UNIQUE,
  address VARCHAR(255) NOT NULL,
  amount DECIMAL(64,20) NOT NULL,
  asset CHAR(32) NOT NULL,
  height INTEGER(11) NOT NULL,
  tx_index INTEGER,
  status INTEGER(11) DEFAULT 1,
  create_time DATETIME DEFAULT NOW(),
  update_time DATETIME DEFAULT NOW() ON UPDATE NOW(),
  PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

# AddScannedDeposit
INSERT INTO deposit (tx_hash, address, amount, asset, height, tx_index) VALUES (?, ?, ?, ?, ?, ?)

# AddDepositWithTime
INSERT INTO deposit (tx_hash, address, amount, asset, height, tx_index, create_time) VALUES (?, ?, ?, ?, ?, ?, ?)

# AddStableDeposit
INSERT INTO deposit (tx_hash, address, amount, asset, height, tx_index, create_time, status) VALUES (?, ?, ?, ?, ?, ?, ?, 2)

# GetUnstableDeposit
SELECT tx_hash, address, amount, asset, height, tx_index FROM deposit WHERE asset=? AND status<2

# DepositIntoStable
UPDATE deposit SET status=2 WHERE tx_hash=?

# GetDepositId
SELECT id FROM deposit WHERE tx_hash=?

# GetDeposits
SELECT id, tx_hash, address, amount, asset, height, tx_index, status, create_time, update_time FROM deposit WHERE %s