CREATE TABLE fines (
    id         SERIAL PRIMARY KEY,
    loan_id    INTEGER NOT NULL UNIQUE REFERENCES loans(id) ON DELETE RESTRICT,
    days_late  INTEGER NOT NULL,
    amount     BIGINT NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'unpaid',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    paid_at    TIMESTAMPTZ,

    CONSTRAINT fines_days_late_positive CHECK (days_late > 0),
    CONSTRAINT fines_amount_positive CHECK (amount > 0),
    CONSTRAINT fines_status_valid CHECK (status IN ('unpaid', 'paid', 'waived')),
    CONSTRAINT fines_paid_at_matches_status CHECK ((status = 'paid') = (paid_at IS NOT NULL))
);