# CreateTable
CREATE TABLE IF NOT EXISTS height (
  asset CHAR(32) NOT NULL UNIQUE,
  height INTEGER(11) DEFAULT 0,
  PRIMARY KEY(asset)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

# ChkExsAsset
SELECT COUNT(asset) FROM height WHERE asset=?

# AddAsset
INSERT INTO height (asset) VALUES (?)

# GetHeight
SELECT height FROM height WHERE asset=?
