CREATE TABLE notifications (
                               id UUID PRIMARY KEY NOT NULL,
                               user_id UUID NOT NULL,
                               title VARCHAR(255) NOT NULL,
                               message TEXT NOT NULL,
                               type VARCHAR(50) NOT NULL DEFAULT 'system',
                               is_read BOOLEAN NOT NULL DEFAULT FALSE,
                               created_at TIMESTAMP NOT NULL,
                               read_at TIMESTAMP NULL
);