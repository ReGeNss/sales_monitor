import { MarketplaceDomain } from '../../marketplaces/domain/marketplace.domain';

export class LatestPriceDomain {
  constructor(
    public readonly priceId: number,
    public readonly regularPrice: number,
    public readonly specialPrice: number | undefined,
    public readonly createdAt: Date,
    public readonly marketplaceProduct: {
      marketplaceProductId: number;
      marketplace: MarketplaceDomain;
      url: string;
      product: {
        productId: number;
        name: string;
        brand: { brandId: number; name: string };
      };
    },
  ) {}
}
