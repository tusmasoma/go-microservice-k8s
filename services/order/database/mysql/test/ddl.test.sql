CREATE DATABASE IF NOT EXISTS `microservice-k8s-demo-test-db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `microservice-k8s-demo-test-db`;

DROP TABLE IF EXISTS OrderLines;
DROP TABLE IF EXISTS Orders;

-- Orders Table
CREATE TABLE Orders (
    id CHAR(36) PRIMARY KEY,
    customer_id CHAR(36) NOT NULL,
    order_date TIMESTAMP NOT NULL
);

-- OrderLines Table
CREATE TABLE OrderLines (
    order_id CHAR(36) NOT NULL,
    catalog_item_id CHAR(36) NOT NULL,
    count INT NOT NULL,
    PRIMARY KEY (order_id, catalog_item_id),
    FOREIGN KEY (order_id) REFERENCES Orders(id)
);
