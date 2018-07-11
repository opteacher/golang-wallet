# CreateTable
CREATE TABLE IF NOT EXISTS withdraw (
  id INTEGER NOT NULL UNIQUE,
  tx_hash VARCHAR(255) UNIQUE,
  address VARCHAR(255) NOT NULL,
  amount DECIMAL(64,20) NOT NULL,
  asset CHAR(32) NOT NULL,
  height INTEGER(11),
  tx_index INTEGER,
  status INTEGER(11) DEFAULT 1,
  create_time DATETIME DEFAULT NOW(),
  update_time DATETIME DEFAULT NOW() ON UPDATE NOW(),
  PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

# RecvNewWithdraw
INSERT INTO withdraw (id, address, amount, asset) VALUES (?, ?, ?, ?)

# SentForTxHash
UPDATE withdraw SET tx_hash=?, status=2 WHERE id=?

# WithdrawIntoStable
UPDATE withdraw SET status=4 WHERE tx_hash=?

# WithdrawIntoChain
UPDATE withdraw SET status=3, %s WHERE tx_hash=?

# GetAvailableId
SELECT MAX(id) + 1 AS new_id FROM withdraw WHERE asset=?

# GetUnstableWithdraw
SELECT id, tx_hash, address, amount, asset, height, tx_index, status, create_time, update_time FROM withdraw WHERE asset=? AND status=3

# GetUnfinishWithdraw
SELECT id, tx_hash, address, amount, asset, height, tx_index, status, create_time, update_time FROM withdraw WHERE asset=? AND status >= 1 AND status < 3

# GetWithdrawId
SELECT id FROM withdraw WHERE tx_hash=?

# GetWithdraws
SELECT id, tx_hash, address, amount, asset, height, tx_index, status, create_time, update_time FROM withdraw WHERE %s