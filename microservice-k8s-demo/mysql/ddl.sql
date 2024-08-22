CREATE DATABASE IF NOT EXISTS `microservice-k8s-demo-db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `microservice-k8s-demo-db`;

DROP TABLE IF EXISTS CatalogItems;
DROP TABLE IF EXISTS Customers;

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