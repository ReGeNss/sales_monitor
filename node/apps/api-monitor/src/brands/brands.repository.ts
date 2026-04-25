import { Injectable } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { Brand } from '@sales-monitor/database';
import { BrandDomain } from './domain/brand.domain';
import { ProductDomain } from '../products/domain/product.domain';
import { CategoryDomain } from '../categories/domain/category.domain';
import { NotFoundError } from '../common/errors';

@Injectable()
export class BrandsRepository {
  constructor(private readonly em: EntityManager) {}

  async findAll(): Promise<BrandDomain[]> {
    const brands = await this.em.find(Brand, {}, { orderBy: { name: 'ASC' } });
    return brands.map((b) => this.toDomain(b));
  }

  async findOne(id: number): Promise<BrandDomain> {
    const brand = await this.em.findOne(Brand, { brandId: id });
    if (!brand) {
      throw new NotFoundError(`Brand with ID ${id} not found`);
    }
    return this.toDomain(brand);
  }

  async findProducts(id: number): Promise<ProductDomain[]> {
    const brand = await this.em.findOne(
      Brand,
      { brandId: id },
      { populate: ['products.category'] },
    );
    if (!brand) {
      throw new NotFoundError(`Brand with ID ${id} not found`);
    }
    return brand.products.getItems().map(
      (p) =>
        new ProductDomain(
          p.productId,
          p.name,
          p.imageUrl,
          this.toDomain(brand),
          new CategoryDomain(p.category.categoryId, p.category.name),
        ),
    );
  }

  private toDomain(orm: Brand): BrandDomain {
    return new BrandDomain(orm.brandId, orm.name, orm.bannerUrl);
  }
}
