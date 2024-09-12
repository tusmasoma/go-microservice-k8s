CREATE DATABASE IF NOT EXISTS `microservice-k8s-demo-db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `microservice-k8s-demo-db`;

GRANT ALL PRIVILEGES ON `microservice-k8s-demo-db`.* TO 'microservice-k8s-demo'@'%';
FLUSH PRIVILEGES;

DROP TABLE IF EXISTS CatalogItems;
DROP TABLE IF EXISTS Customers;
DROP TABLE IF EXISTS OrderLines;
DROP TABLE IF EXISTS Orders;

-- CatalogItems Table
CREATE TABLE CatalogItems (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

-- Customers Table
CREATE TABLE Customers (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    street VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL
);

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
