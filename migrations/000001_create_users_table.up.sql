    Create Table users(                              
      id  Serial PRIMARY key,
      name VARCHAR(50) NOT NULL,
      email VARCHAR(100) UNIQUE NOT NULL,
      password_hash VARCHAR(255) NOT NULL,
      role VARCHAR(50) DEFAULT 'member'

 CONSTRAINT users_name_not_blank CHECK (length(trim(name)) > 0),
    CONSTRAINT users_email_lowercase CHECK (email = lower(email)),
    CONSTRAINT users_role_valid CHECK (role IN ('member', 'admin'))
); 