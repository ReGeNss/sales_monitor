import { Injectable } from '@nestjs/common';
import { PricesRepository } from './prices.repository';
import { LatestPriceDomain } from './domain/latest-price.domain';

@Injectable()
export class PricesService {
  constructor(private readonly pricesRepository: PricesRepository) {}

  async getLatestPrices(limit = 100): Promise<LatestPriceDomain[]> {
    return this.pricesRepository.findLatest(limit);
  }

  async getPriceTrends(productId?: number, days = 30): Promise<LatestPriceDomain[]> {
    return this.pricesRepository.findTrends(productId, days);
  }
}
