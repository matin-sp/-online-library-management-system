CREATE TABLE users(
    Create Table User(                              
      id  Serial PRIMARY key,
      name VARCHAR(50) NOT NULL,
      email VARCHAR(100) UNIQUE NOT NULL,
      password_hash VARCHAR(255) NOT NULL,
      role VARCHAR(50) DEFAULT 'member',
); 