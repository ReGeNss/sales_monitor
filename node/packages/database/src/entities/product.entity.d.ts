import { Collection } from '@mikro-orm/core';
import { Brand } from './brand.entity';
import { Category } from './category.entity';
import { MarketplaceProduct } from './marketplace-product.entity';
import { ProductAttribute } from './product-attribute.entity';
export declare class Product {
    productId: number;
    nameFingerprint?: string;
    brandId: number;
    brand: Brand;
    name: string;
    categoryId: number;
    category: Category;
    imageUrl?: string;
    marketplaceProducts: Collection<MarketplaceProduct, object>;
    attributes: Collection<ProductAttribute, object>;
}
