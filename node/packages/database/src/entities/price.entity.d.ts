import { MarketplaceProduct } from './marketplace-product.entity';
export declare class Price {
    priceId: number;
    marketplaceProductId: number;
    marketplaceProduct: MarketplaceProduct;
    regularPrice: number;
    discountPrice?: number;
    createdAt: Date;
}
