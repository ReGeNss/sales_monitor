-- Modify "products" table
ALTER TABLE `products` ADD COLUMN `name_fingerprint` varchar(191) NOT NULL AFTER `name`, ADD INDEX `idx_product_fingerprint_fulltext` (`name_fingerprint`), ADD UNIQUE INDEX `uni_products_name_fingerprint` (`name_fingerprint`);
