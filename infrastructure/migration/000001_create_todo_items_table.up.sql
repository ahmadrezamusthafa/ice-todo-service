CREATE TABLE IF NOT EXISTS todo_items (
    id CHAR(36) PRIMARY KEY,
    description TEXT NOT NULL,
    due_date TIMESTAMP NOT NULL,
    file_id VARCHAR(255)
);