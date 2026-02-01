import { Injectable, NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { Category } from '@sales-monitor/database';

@Injectable()
export class CategoriesService {
  constructor(private readonly em: EntityManager) {}

  async findAll() {
    return this.em.find(Category, {}, { orderBy: { name: 'ASC' } });
  }

  async findOne(id: number) {
    const category = await this.em.findOne(Category, { categoryId: id });
    if (!category) {
      throw new NotFoundException(`Category with ID ${id} not found`);
    }
    return category;
  }

  async getCategoryProducts(categoryId: number, page = 1, limit = 20) {
    const category = await this.em.findOne(
      Category,
      { categoryId },
      { populate: ['products.brand'] },
    );

    if (!category) {
      throw new NotFoundException(`Category with ID ${categoryId} not found`);
    }

    const products = category.products.getItems();
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
