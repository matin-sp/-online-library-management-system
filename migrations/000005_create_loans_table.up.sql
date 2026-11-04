CREATE TABLE loans (
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    copy_id       INTEGER NOT NULL REFERENCES book_copies(id) ON DELETE RESTRICT,
    borrowed_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    due_date      TIMESTAMPTZ NOT NULL,
    returned_at   TIMESTAMPTZ,
    renewal_count INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT loans_due_after_borrowed CHECK (due_date > borrowed_at),
    CONSTRAINT loans_returned_after_borrowed
        CHECK (returned_at IS NULL OR returned_at >= borrowed_at),
    CONSTRAINT loans_renewal_count_valid CHECK (renewal_count BETWEEN 0 AND 1)
);

CREATE UNIQUE INDEX uq_active_loan
    ON loans (copy_id)
    WHERE returned_at IS NULL;

CREATE INDEX loans_user_history_idx ON loans (user_id, borrowed_at DESC);

CREATE INDEX loans_due_date_open_idx
    ON loans (due_date)
    WHERE returned_at IS NULL;