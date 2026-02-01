import { Collection } from '@mikro-orm/core';
import { Marketplace } from './marketplace.entity';
import { Product } from './product.entity';
import { Price } from './price.entity';
export declare class MarketplaceProduct {
    marketplaceProductId: number;
    marketplaceId: number;
    marketplace: Marketplace;
    productId: number;
    product: Product;
    url: string;
    prices: Collection<Price, object>;
}
