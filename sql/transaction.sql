# CreateTable
CREATE TABLE IF NOT EXISTS `transaction` (
  opr_info VARCHAR(255) UNIQUE COMMENT 'type_id',
  tx_hash VARCHAR(255) NOT NULL,
  block_hash VARCHAR(255) NOT NULL,
  `from` VARCHAR(255) NOT NULL,
  `to` VARCHAR(255) NOT NULL,
  amount DECIMAL(64,20) NOT NULL,
  asset CHAR(32) NOT NULL,
  height INTEGER(11) DEFAULT 0,
  tx_index INTEGER DEFAULT 0,
  create_time DATETIME DEFAULT NOW(),
  PRIMARY KEY (opr_info)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

# AddTransaction
INSERT INTO `transaction` (%s) VALUES (%s)