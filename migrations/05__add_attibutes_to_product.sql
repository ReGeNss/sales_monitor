CREATE TABLE product_attributes(
    product_attribute_id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL REFERENCES Product(product_id),
    attribute_type ENUM('volume', 'weight') NOT NULL,
    value TEXT NOT NULL
);