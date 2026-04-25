import { MarketplaceDomain } from '../../marketplaces/domain/marketplace.domain';
import { PriceDomain } from '../../prices/domain/price.domain';

export class ProductPriceDomain {
  constructor(
    public readonly marketplace: MarketplaceDomain,
    public readonly url: string,
    public readonly prices: PriceDomain[],
  ) {}
}
