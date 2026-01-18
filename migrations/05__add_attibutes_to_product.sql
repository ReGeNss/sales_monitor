CREATE TABLE attributes(
    attribute_id INT AUTO_INCREMENT PRIMARY KEY,
    attribute_type ENUM('volume', 'weight') NOT NULL,
    value TEXT NOT NULL
);