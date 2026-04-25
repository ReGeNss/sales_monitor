import { Injectable } from '@nestjs/common';
import { createPaginationMeta } from '@sales-monitor/common';
import { CategoriesRepository } from './categories.repository';
import { CategoryDomain } from './domain/category.domain';
import { ProductDomain } from '../products/domain/product.domain';

@Injectable()
export class CategoriesService {
  constructor(private readonly categoriesRepository: CategoriesRepository) {}

  async findAll(): Promise<CategoryDomain[]> {
    return this.categoriesRepository.findAll();
  }

  async findOne(id: number): Promise<CategoryDomain> {
    return this.categoriesRepository.findOne(id);
  }

  async getCategoryProducts(categoryId: number, page = 1, limit = 20) {
    const products = await this.categoriesRepository.findProducts(categoryId);
    const total = products.length;
    const data: ProductDomain[] = products.slice((page - 1) * limit, page * limit);

    return {
      data,
      meta: createPaginationMeta(page, limit, total),
    };
  }
}
