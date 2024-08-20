CREATE DATABASE IF NOT EXISTS `microservice-k8s-demo-db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `microservice-k8s-demo-db`;

DROP TABLE IF EXISTS CatalogItems;

-- CatalogItems Table
CREATE TABLE CatalogItems (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);