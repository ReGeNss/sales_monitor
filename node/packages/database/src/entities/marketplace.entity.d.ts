import { Collection } from '@mikro-orm/core';
import { MarketplaceProduct } from './marketplace-product.entity';
export declare class Marketplace {
    marketplaceId: number;
    name: string;
    url: string;
    marketplaceProducts: Collection<MarketplaceProduct, object>;
}
