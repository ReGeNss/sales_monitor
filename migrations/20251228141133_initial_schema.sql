-- Create "brands" table
CREATE TABLE `brands` (
  `brand_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext NOT NULL,
  `banner_url` longtext NULL,
  PRIMARY KEY (`brand_id`)
) COLLATE utf8mb4_uca1400_ai_ci;
-- Create "categories" table
CREATE TABLE `categories` (
  `category_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext NOT NULL,
  `fieldname` longtext NULL,
  PRIMARY KEY (`category_id`)
) COLLATE utf8mb4_uca1400_ai_ci;
-- Create "favorite_brand" table
CREATE TABLE `favorite_brand` (
  `user_id` bigint unsigned NOT NULL,
  `brand_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`, `brand_id`),
  INDEX `fk_favorite_brand_brand` (`brand_id`)
) COLLATE utf8mb4_uca1400_ai_ci;
-- Create "favorite_brands" table
CREATE TABLE `favorite_brands` (
  `favorite_brand_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NULL,
  `brand_id` bigint unsigned NULL,
  PRIMARY KEY (`favorite_brand_id`)
) COLLATE utf8mb4_uca1400_ai_ci;
-- Create "favorite_product" table
CREATE TABLE `favorite_product` (
  `user_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`, `product_id`),
  INDEX `fk_favorite_product_product` (`product_id`)
) COLLATE utf8mb4_uca1400_ai_ci;
-- Create "favorite_products" table
CREATE TABLE `favorite_products` (
  `favorite_product_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NULL,
  `product_id` bigint unsigned NULL,
  PRIMARY KEY (`favorite_product_id`)
) COLLATE utf8mb4_uca1400_ai_ci;
-- Create "marketplaces" table
CREATE TABLE `marketplaces` (
  `marketplace_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext NOT NULL,
  `url` longtext NOT NULL,
  PRIMARY KEY (`marketplace_id`)
) COLLATE utf8mb4_uca1400_ai_ci;
-- Create "prices" table
CREATE TABLE `prices` (
  `price_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `product_id` bigint unsigned NULL,
  `marketplace_id` bigint unsigned NULL,
  `regular_price` decimal(10,2) NOT NULL,
  `discount_price` decimal(10,2) NULL,
  `url` longtext NOT NULL,
  `is_on_sale` bool NULL DEFAULT 0,
  `discount_percent` bigint NULL,
  `created_at` datetime(3) NULL,
  PRIMARY KEY (`price_id`),
  INDEX `fk_products_prices` (`product_id`)
) COLLATE utf8mb4_uca1400_ai_ci;
-- Create "products" table
CREATE TABLE `products` (
  `product_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext NOT NULL,
  `image_url` longtext NULL,
  `brand_id` bigint unsigned NULL,
  `category_id` bigint unsigned NULL,
  PRIMARY KEY (`product_id`),
  INDEX `fk_brands_products` (`brand_id`),
  INDEX `fk_categories_products` (`category_id`)
) COLLATE utf8mb4_uca1400_ai_ci;
-- Create "users" table
CREATE TABLE `users` (
  `user_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `login` varchar(191) NOT NULL,
  `password` longtext NOT NULL,
  `nf_token` longtext NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE INDEX `uni_users_login` (`login`)
) COLLATE utf8mb4_uca1400_ai_ci;
-- Modify "brands" table
ALTER TABLE `brands` ADD CONSTRAINT `fk_products_brand` FOREIGN KEY (`brand_id`) REFERENCES `products` (`product_id`) ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "categories" table
ALTER TABLE `categories` ADD CONSTRAINT `fk_products_category` FOREIGN KEY (`category_id`) REFERENCES `products` (`product_id`) ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "favorite_brand" table
ALTER TABLE `favorite_brand` ADD CONSTRAINT `fk_favorite_brand_brand` FOREIGN KEY (`brand_id`) REFERENCES `brands` (`brand_id`) ON UPDATE RESTRICT ON DELETE RESTRICT, ADD CONSTRAINT `fk_favorite_brand_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "favorite_product" table
ALTER TABLE `favorite_product` ADD CONSTRAINT `fk_favorite_product_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`product_id`) ON UPDATE RESTRICT ON DELETE RESTRICT, ADD CONSTRAINT `fk_favorite_product_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "marketplaces" table
ALTER TABLE `marketplaces` ADD CONSTRAINT `fk_prices_marketplace` FOREIGN KEY (`marketplace_id`) REFERENCES `prices` (`price_id`) ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "prices" table
ALTER TABLE `prices` ADD CONSTRAINT `fk_products_prices` FOREIGN KEY (`product_id`) REFERENCES `products` (`product_id`) ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "products" table
ALTER TABLE `products` ADD CONSTRAINT `fk_brands_products` FOREIGN KEY (`brand_id`) REFERENCES `brands` (`brand_id`) ON UPDATE RESTRICT ON DELETE RESTRICT, ADD CONSTRAINT `fk_categories_products` FOREIGN KEY (`category_id`) REFERENCES `categories` (`category_id`) ON UPDATE RESTRICT ON DELETE RESTRICT, ADD CONSTRAINT `fk_prices_product` FOREIGN KEY (`product_id`) REFERENCES `prices` (`price_id`) ON UPDATE RESTRICT ON DELETE RESTRICT;
