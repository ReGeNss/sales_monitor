import { Injectable, NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { Product } from '@sales-monitor/database';
import { createPaginationMeta, PaginatedResponse } from '@sales-monitor/common';
import { ProductFilterDto } from './dto/product-filter.dto';

@Injectable()
export class ProductsService {
  constructor(private readonly em: EntityManager) {}

  async findAll(filterDto: ProductFilterDto): Promise<PaginatedResponse<Product>> {
    const { page = 1, limit = 20, categoryId, brandId, search } = filterDto;
    const where: any = {};

    if (categoryId) {
      where.category = { categoryId };
    }

    if (brandId) {
      where.brand = { brandId };
    }

    if (search) {
      where.name = { $like: `%${search}%` };
    }

    const [products, total] = await this.em.findAndCount(
      Product,
      where,
      {
        populate: ['brand', 'category'],
        limit,
        offset: (page - 1) * limit,
        orderBy: { productId: 'DESC' },
      },
    );

    return {
      data: products,
      meta: createPaginationMeta(page, limit, total),
    };
  }

  async findOne(id: number): Promise<Product> {
    const product = await this.em.findOne(
      Product,
      { productId: id },
      { populate: ['brand', 'category', 'attributes', 'marketplaceProducts.marketplace'] },
    );

    if (!product) {
      throw new NotFoundException(`Product with ID ${id} not found`);
    }

    return product;
  }

  async getProductPrices(productId: number) {
    const product = await this.em.findOne(
      Product,
      { productId },
      { populate: ['marketplaceProducts.prices', 'marketplaceProducts.marketplace'] },
    );

    if (!product) {
      throw new NotFoundException(`Product with ID ${productId} not found`);
    }

    return product.marketplaceProducts.getItems().map((mp) => ({
      marketplace: mp.marketplace,
      url: mp.url,
      prices: mp.prices.getItems().sort((a, b) => 
        b.createdAt.getTime() - a.createdAt.getTime()
      ).slice(0, 30), // Last 30 price records per marketplace
    }));
  }
}
