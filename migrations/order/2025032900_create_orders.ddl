CREATE TABLE Orders (
    id CHAR(36) PRIMARY KEY,
    customer_id CHAR(36) NOT NULL,
    order_date TIMESTAMP NOT NULL
);

CREATE TABLE OrderLines (
    order_id CHAR(36) NOT NULL,
    catalog_item_id CHAR(36) NOT NULL,
    count INT NOT NULL,
    PRIMARY KEY (order_id, catalog_item_id),
    FOREIGN KEY (order_id) REFERENCES Orders(id)
);