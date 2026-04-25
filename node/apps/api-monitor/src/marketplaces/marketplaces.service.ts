import { Injectable } from '@nestjs/common';
import { MarketplacesRepository } from './marketplaces.repository';
import { MarketplaceDomain } from './domain/marketplace.domain';

@Injectable()
export class MarketplacesService {
  constructor(private readonly marketplacesRepository: MarketplacesRepository) {}

  async findAll(): Promise<MarketplaceDomain[]> {
    return this.marketplacesRepository.findAll();
  }

  async findOne(id: number): Promise<MarketplaceDomain> {
    return this.marketplacesRepository.findOne(id);
  }
}
