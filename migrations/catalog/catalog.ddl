CREATE TABLE CatalogItems (
  id char(36) NOT NULL,
  name varchar(255) NOT NULL,
  price decimal(10,2) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_catalogitems_name (name)
)
