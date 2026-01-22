CREATE TABLE product_attributes(
    product_id INT NOT NULL,
    attribute_id INT NOT NULL,
    PRIMARY KEY (product_id, attribute_id),
    CONSTRAINT fk_pa_product FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    CONSTRAINT fk_pa_attribute FOREIGN KEY (attribute_id) REFERENCES attributes(attribute_id) ON DELETE CASCADE
);

CREATE INDEX idx_pa_product_id ON product_attributes(product_id);
CREATE INDEX idx_pa_attribute_id ON product_attributes(attribute_id);
