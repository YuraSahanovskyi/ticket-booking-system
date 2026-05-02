CREATE TABLE seats (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    row INT NOT NULL,
    number INT NOT NULL,

    UNIQUE(event_id, row, number)
);