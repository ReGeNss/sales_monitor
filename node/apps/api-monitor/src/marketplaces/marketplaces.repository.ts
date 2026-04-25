import { Injectable } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { Marketplace } from '@sales-monitor/database';
import { MarketplaceDomain } from './domain/marketplace.domain';
import { NotFoundError } from '../common/errors';

@Injectable()
export class MarketplacesRepository {
  constructor(private readonly em: EntityManager) {}

  async findAll(): Promise<MarketplaceDomain[]> {
    const marketplaces = await this.em.find(Marketplace, {}, { orderBy: { name: 'ASC' } });
    return marketplaces.map((m) => this.toDomain(m));
  }

  async findOne(id: number): Promise<MarketplaceDomain> {
    const marketplace = await this.em.findOne(Marketplace, { marketplaceId: id });
    if (!marketplace) {
      throw new NotFoundError(`Marketplace with ID ${id} not found`);
    }
    return this.toDomain(marketplace);
  }

  private toDomain(orm: Marketplace): MarketplaceDomain {
    return new MarketplaceDomain(orm.marketplaceId, orm.name, orm.url);
  }
}
