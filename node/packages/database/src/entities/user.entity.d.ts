import { Collection } from '@mikro-orm/core';
import { Product } from './product.entity';
import { Brand } from './brand.entity';
export declare class User {
    userId: number;
    login: string;
    password: string;
    nfToken?: string;
    favoriteProducts: Collection<Product, object>;
    favoriteBrands: Collection<Brand, object>;
}
