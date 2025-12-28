-- Modify "brands" table
ALTER TABLE `brands` MODIFY COLUMN `name` varchar(191) NOT NULL, ADD UNIQUE INDEX `uni_brands_name` (`name`);
-- Modify "categories" table
ALTER TABLE `categories` MODIFY COLUMN `name` varchar(191) NOT NULL, DROP COLUMN `fieldname`, ADD UNIQUE INDEX `uni_categories_name` (`name`);
-- Modify "favorite_brands" table
ALTER TABLE `favorite_brands` MODIFY COLUMN `user_id` bigint unsigned NOT NULL, MODIFY COLUMN `brand_id` bigint unsigned NOT NULL;
-- Modify "favorite_products" table
ALTER TABLE `favorite_products` MODIFY COLUMN `user_id` bigint unsigned NOT NULL, MODIFY COLUMN `product_id` bigint unsigned NOT NULL;
-- Modify "marketplaces" table
ALTER TABLE `marketplaces` MODIFY COLUMN `name` varchar(191) NOT NULL, MODIFY COLUMN `url` varchar(191) NOT NULL, ADD UNIQUE INDEX `uni_marketplaces_name` (`name`), ADD UNIQUE INDEX `uni_marketplaces_url` (`url`);
-- Modify "prices" table
ALTER TABLE `prices` MODIFY COLUMN `product_id` bigint unsigned NOT NULL, MODIFY COLUMN `marketplace_id` bigint unsigned NOT NULL, MODIFY COLUMN `is_on_sale` bool NOT NULL DEFAULT 0, MODIFY COLUMN `discount_percent` bigint NOT NULL DEFAULT 0;
-- Modify "products" table
ALTER TABLE `products` MODIFY COLUMN `name` varchar(191) NOT NULL, MODIFY COLUMN `brand_id` bigint unsigned NOT NULL, MODIFY COLUMN `category_id` bigint unsigned NOT NULL, ADD UNIQUE INDEX `uni_products_name` (`name`);
