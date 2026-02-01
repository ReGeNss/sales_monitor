import { Collection } from '@mikro-orm/core';
import { Product } from './product.entity';
export declare class Brand {
    brandId: number;
    name: string;
    bannerUrl?: string;
    products: Collection<Product, object>;
}
