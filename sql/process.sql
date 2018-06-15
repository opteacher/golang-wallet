# CreateTable
CREATE TABLE IF NOT EXISTS process (
  tx_hash VARCHAR(255) NOT NULL UNIQUE,
  `type` VARCHAR(64) NOT NULL COMMENT 'DEPOSIT/COLLECT/WITHDRAW',
  height INTEGER(11) DEFAULT 0,
  complete_height INTEGER(11),
  process VARCHAR(64) NOT NULL,
  cancelable TINYINT(1) NOT NULL DEFAULT 1 COMMENT '0/1: 不/可取消',
  last_update_time DATETIME DEFAULT NOW() ON UPDATE NOW(),
  PRIMARY KEY(tx_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

# AddProcess
INSERT INTO process (tx_hash, `type`, height, complete_height, process, cancelable) VALUES (?, ?, ?, ?, ?, ?)

# UpdateProcessByHash
UPDATE process SET %s WHERE tx_hash=?

# CheckProcsExists
SELECT COUNT(tx_hash) FROM process WHERE tx_hash=?