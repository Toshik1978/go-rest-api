CREATE TABLE payments(
                         id BIGSERIAL PRIMARY KEY,
                         amount INT NOT NULL,
                         payer_account_uid VARCHAR(256) NOT NULL,
                         recipient_account_uid VARCHAR(256) NOT NULL,
                         created_at TIMESTAMP WITH TIME ZONE NOT NULL,
                         FOREIGN KEY (payer_account_uid) REFERENCES accounts(uid),
                         FOREIGN KEY (recipient_account_uid) REFERENCES accounts(uid)
);

CREATE INDEX ON payments(payer_account_uid);
CREATE INDEX ON payments(recipient_account_uid);
