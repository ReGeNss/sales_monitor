import { Collection } from '@mikro-orm/core';
import { Product } from './product.entity';
export declare class Category {
    categoryId: number;
    name: string;
    products: Collection<Product, object>;
}
