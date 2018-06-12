# CreateTable
CREATE TABLE IF NOT EXISTS address (
  id INTEGER NOT NULL AUTO_INCREMENT,
  asset VARCHAR(255) NOT NULL,
  address VARCHAR(255) NOT NULL UNIQUE,
  inuse TINYINT(1) NOT NULL DEFAULT 0,
  PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

# NewAddress
INSERT INTO address (asset, address) VALUES (?, ?)

# FindByAsset
SELECT address FROM address WHERE inuse=1 AND asset=?