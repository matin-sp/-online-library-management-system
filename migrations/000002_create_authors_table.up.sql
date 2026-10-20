CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,

    CONSTRAINT authors_name_not_blank CHECK (length(trim(name)) > 0)
);

CREATE UNIQUE INDEX authors_name_unique ON authors (lower(name));