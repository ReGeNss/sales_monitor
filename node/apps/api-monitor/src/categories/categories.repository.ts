import { Injectable } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { Category } from '@sales-monitor/database';
import { CategoryDomain } from './domain/category.domain';
import { ProductDomain } from '../products/domain/product.domain';
import { BrandDomain } from '../brands/domain/brand.domain';
import { NotFoundError } from '../common/errors';

@Injectable()
export class CategoriesRepository {
  constructor(private readonly em: EntityManager) {}

  async findAll(): Promise<CategoryDomain[]> {
    const categories = await this.em.find(Category, {}, { orderBy: { name: 'ASC' } });
    return categories.map((c) => this.toDomain(c));
  }

  async findOne(id: number): Promise<CategoryDomain> {
    const category = await this.em.findOne(Category, { categoryId: id });
    if (!category) {
      throw new NotFoundError(`Category with ID ${id} not found`);
    }
    return this.toDomain(category);
  }

  async findProducts(id: number): Promise<ProductDomain[]> {
    const category = await this.em.findOne(
      Category,
      { categoryId: id },
      { populate: ['products.brand'] },
    );
    if (!category) {
      throw new NotFoundError(`Category with ID ${id} not found`);
    }
    return category.products.getItems().map(
      (p) =>
        new ProductDomain(
          p.productId,
          p.name,
          p.imageUrl,
          new BrandDomain(p.brand.brandId, p.brand.name, p.brand.bannerUrl),
          this.toDomain(category),
        ),
    );
  }

  private toDomain(orm: Category): CategoryDomain {
    return new CategoryDomain(orm.categoryId, orm.name);
  }
}
