import { Injectable } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { Price } from '@sales-monitor/database';

@Injectable()
export class PricesService {
  constructor(private readonly em: EntityManager) {}

  async getLatestPrices(limit = 100) {
    const prices = await this.em.find(
      Price,
      {},
      {
        populate: ['marketplaceProduct.product.brand', 'marketplaceProduct.marketplace'],
        orderBy: { createdAt: 'DESC' },
        limit,
      },
    );

    return prices;
  }

  async getPriceTrends(productId?: number, days = 30) {
    const dateFrom = new Date();
    dateFrom.setDate(dateFrom.getDate() - days);

    const where: any = {
      createdAt: { $gte: dateFrom },
    };

    if (productId) {
      where['marketplaceProduct.product.productId'] = productId;
    }

    const prices = await this.em.find(
      Price,
      where,
      {
        populate: ['marketplaceProduct.product', 'marketplaceProduct.marketplace'],
        orderBy: { createdAt: 'ASC' },
      },
    );

    return prices;
  }
}
