CREATE TABLE book_copies (
    id          SERIAL PRIMARY KEY,
    book_id     INTEGER NOT NULL REFERENCES books(id) ON DELETE RESTRICT,
    copy_number INTEGER NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'available',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT book_copies_copy_number_positive CHECK (copy_number > 0),
    CONSTRAINT book_copies_status_valid
        CHECK (status IN ('available', 'checked_out', 'under_repair', 'lost', 'damaged')),
    CONSTRAINT book_copies_book_copy_unique UNIQUE (book_id, copy_number)
);