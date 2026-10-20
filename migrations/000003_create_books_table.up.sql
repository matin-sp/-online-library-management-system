CREATE TABLE books (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    release_year INTEGER,
    author_id    INTEGER NOT NULL REFERENCES authors(id) ON DELETE RESTRICT,

    CONSTRAINT books_title_not_blank CHECK (length(trim(title)) > 0),
    CONSTRAINT books_year_valid CHECK (release_year IS NULL OR release_year BETWEEN 1000 AND 2100),
    CONSTRAINT books_title_author_unique UNIQUE (title, author_id)
);

CREATE INDEX books_author_id_idx ON books (author_id);
