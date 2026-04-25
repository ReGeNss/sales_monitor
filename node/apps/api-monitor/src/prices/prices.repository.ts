import { Injectable } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { Price } from '@sales-monitor/database';
import { LatestPriceDomain } from './domain/latest-price.domain';
import { MarketplaceDomain } from '../marketplaces/domain/marketplace.domain';

@Injectable()
export class PricesRepository {
  constructor(private readonly em: EntityManager) {}

  async findLatest(limit = 100): Promise<LatestPriceDomain[]> {
    const prices = await this.em.find(
      Price,
      {},
      {
        populate: ['marketplaceProduct.product.brand', 'marketplaceProduct.marketplace'],
        orderBy: { createdAt: 'DESC' },
        limit,
      },
    );
    return prices.map((p) => this.toLatestDomain(p));
  }

  async findTrends(productId?: number, days = 30): Promise<LatestPriceDomain[]> {
    const dateFrom = new Date();
    dateFrom.setDate(dateFrom.getDate() - days);

    const where: any = { createdAt: { $gte: dateFrom } };
    if (productId) {
      where['marketplaceProduct.product.productId'] = productId;
    }

    const prices = await this.em.find(Price, where, {
      populate: ['marketplaceProduct.product', 'marketplaceProduct.marketplace'],
      orderBy: { createdAt: 'ASC' },
    });
    return prices.map((p) => this.toLatestDomain(p));
  }

  private toLatestDomain(orm: Price): LatestPriceDomain {
    const mp = orm.marketplaceProduct;
    return new LatestPriceDomain(
      orm.priceId,
      orm.regularPrice,
      orm.specialPrice,
      orm.createdAt,
      {
        marketplaceProductId: mp.marketplaceProductId,
        marketplace: new MarketplaceDomain(mp.marketplace.marketplaceId, mp.marketplace.name, mp.marketplace.url),
        url: mp.url,
        product: {
          productId: mp.product.productId,
          name: mp.product.name,
          brand: { brandId: mp.product.brand.brandId, name: mp.product.brand.name },
        },
      },
    );
  }
}
