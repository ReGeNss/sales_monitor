import { MarketplaceDomain } from '../../marketplaces/domain/marketplace.domain';
import { PriceDomain } from '../../prices/domain/price.domain';

export class MarketplaceProductDomain {
  constructor(
    public readonly marketplaceProductId: number,
    public readonly marketplace: MarketplaceDomain,
    public readonly url: string,
    public readonly prices: PriceDomain[],
  ) {}
}
