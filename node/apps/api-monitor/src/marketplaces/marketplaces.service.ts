import { Injectable, NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { Marketplace } from '@sales-monitor/database';

@Injectable()
export class MarketplacesService {
  constructor(private readonly em: EntityManager) {}

  async findAll() {
    return this.em.find(Marketplace, {}, { orderBy: { name: 'ASC' } });
  }

  async findOne(id: number) {
    const marketplace = await this.em.findOne(Marketplace, { marketplaceId: id });
    if (!marketplace) {
      throw new NotFoundException(`Marketplace with ID ${id} not found`);
    }
    return marketplace;
  }
}
