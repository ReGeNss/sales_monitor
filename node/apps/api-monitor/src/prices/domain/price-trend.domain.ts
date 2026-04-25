import { MarketplaceDomain } from '../../marketplaces/domain/marketplace.domain';
import { PriceDomain } from './price.domain';

export class PriceTrendDomain {
  constructor(
    public readonly marketplace: MarketplaceDomain,
    public readonly productId: number,
    public readonly prices: PriceDomain[],
  ) {}
}
