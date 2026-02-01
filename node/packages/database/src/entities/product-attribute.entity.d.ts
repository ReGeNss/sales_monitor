import { Collection } from '@mikro-orm/core';
import { Product } from './product.entity';
export declare class ProductAttribute {
    attributeId: number;
    name: string;
    value: string;
    products: Collection<Product, object>;
}
