CREATE TABLE OrderLines (
  order_id char(36) NOT NULL,
  catalog_item_id char(36) NOT NULL,
  quantity int NOT NULL,
  PRIMARY KEY (order_id,catalog_item_id),
  CONSTRAINT OrderLines_ibfk_1 FOREIGN KEY (order_id) REFERENCES Orders (id)
)
CREATE TABLE Orders (
  id char(36) NOT NULL,
  customer_id char(36) NOT NULL,
  order_date timestamp NOT NULL,
  PRIMARY KEY (id)
)
