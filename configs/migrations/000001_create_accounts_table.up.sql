CREATE TABLE accounts(
                         id BIGSERIAL PRIMARY KEY,
                         uid VARCHAR(256) NOT NULL,
                         currency VARCHAR(16) NOT NULL,
                         balance INT NOT NULL CHECK (balance > 0),
                         created_at TIMESTAMP WITH TIME ZONE NOT NULL,
                         UNIQUE(uid)
);
