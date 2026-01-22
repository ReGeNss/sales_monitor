CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nf_token TEXT
);

CREATE TABLE categories (
    category_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE brands (
    brand_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    banner_url TEXT
);

CREATE TABLE marketplaces (
    marketplace_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    url TEXT NOT NULL
);


CREATE TABLE products (
    product_id INT AUTO_INCREMENT PRIMARY KEY,
    name_fingerprint VARCHAR(255) CHECK (name_fingerprint <> ''),
    brand_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    category_id INT NOT NULL,
    image_url TEXT,
    CONSTRAINT fk_product_brand FOREIGN KEY (brand_id) REFERENCES brands(brand_id),
    CONSTRAINT fk_product_category FOREIGN KEY (category_id) REFERENCES categories(category_id)
);

CREATE TABLE marketplace_products (
    marketplace_product_id INT AUTO_INCREMENT PRIMARY KEY,
    marketplace_id INT NOT NULL,
    product_id INT NOT NULL,
    url TEXT NOT NULL,
    CONSTRAINT fk_marketplace_product_marketplace FOREIGN KEY (marketplace_id) REFERENCES marketplaces(marketplace_id),
    CONSTRAINT fk_marketplace_product_product FOREIGN KEY (product_id) REFERENCES products(product_id)
);

CREATE TABLE prices (
    price_id INT AUTO_INCREMENT PRIMARY KEY,
    marketplace_product_id INT NOT NULL,
    regular_price DECIMAL(10, 2) NOT NULL,
    discount_price DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_price_marketplace_product FOREIGN KEY (marketplace_product_id) REFERENCES marketplace_products(marketplace_product_id)
);

CREATE TABLE favorite_products (
    favorite_product_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    CONSTRAINT fk_fav_prod_user FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_fav_prod_product FOREIGN KEY (product_id) REFERENCES products(product_id)
);

CREATE TABLE favorite_brands (
    favorite_brand_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    brand_id INT NOT NULL,
    CONSTRAINT fk_fav_brand_user FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_fav_brand_brand FOREIGN KEY (brand_id) REFERENCES brands(brand_id)
);