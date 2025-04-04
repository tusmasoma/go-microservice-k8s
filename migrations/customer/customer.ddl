CREATE TABLE Customers (
  id char(36) NOT NULL,
  name varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  street varchar(255) NOT NULL,
  city varchar(255) NOT NULL,
  country_code char(3) NOT NULL DEFAULT 'JPN',
  PRIMARY KEY (id)
)
