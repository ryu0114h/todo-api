CREATE TABLE tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    due_date DATETIME,
    assignee_id INT,
    create_user_id INT NOT NULL,
    visibility VARCHAR(10) NOT NULL DEFAULT 'private',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (assignee_id) REFERENCES users (id),
    FOREIGN KEY (create_user_id) REFERENCES users (id),
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);