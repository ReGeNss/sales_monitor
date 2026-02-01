import { Injectable, NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { Brand } from '@sales-monitor/database';

@Injectable()
export class BrandsService {
  constructor(private readonly em: EntityManager) {}

  async findAll() {
    return this.em.find(Brand, {}, { orderBy: { name: 'ASC' } });
  }

  async findOne(id: number) {
    const brand = await this.em.findOne(Brand, { brandId: id });
    if (!brand) {
      throw new NotFoundException(`Brand with ID ${id} not found`);
    }
    return brand;
  }

  async getBrandProducts(brandId: number, page = 1, limit = 20) {
    const brand = await this.em.findOne(
      Brand,
      { brandId },
      { populate: ['products.category'] },
    );

    if (!brand) {
      throw new NotFoundException(`Brand with ID ${brandId} not found`);
    }

    const products = brand.products.getItems();
    const total = products.length;
    const paginatedProducts = products.slice((page - 1) * limit, page * limit);

    return {
      data: paginatedProducts,
      meta: {
        page,
        limit,
        total,
        totalPages: Math.ceil(total / limit),
      },
    };
  }
}
