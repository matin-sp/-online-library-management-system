CREATE TABLE books (
    id SERIAL Primary key,
    title VARCHAR(255) NOT NULL,
    release_year INTEGER,
    author_id INTEGER REFERENCES authors(id)
);



