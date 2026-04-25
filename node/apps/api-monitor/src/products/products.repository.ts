import { Injectable } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { Product } from '@sales-monitor/database';
import { createPaginationMeta, PaginatedResponse } from '@sales-monitor/common';
import { ProductDomain } from './domain/product.domain';
import { ProductDetailDomain } from './domain/product-detail.domain';
import { ProductPriceDomain } from './domain/product-price.domain';
import { ProductAttributeDomain } from './domain/product-attribute.domain';
import { MarketplaceProductDomain } from './domain/marketplace-product.domain';
import { BrandDomain } from '../brands/domain/brand.domain';
import { CategoryDomain } from '../categories/domain/category.domain';
import { MarketplaceDomain } from '../marketplaces/domain/marketplace.domain';
import { PriceDomain } from '../prices/domain/price.domain';
import { NotFoundError } from '../common/errors';
import { ProductFilterDto } from './dto/product-filter.dto';

@Injectable()
export class ProductsRepository {
  constructor(private readonly em: EntityManager) {}

  async findAll(filterDto: ProductFilterDto): Promise<PaginatedResponse<ProductDomain>> {
    const { page = 1, limit = 20, categoryId, brandId, search } = filterDto;
    const where: any = {};

    if (categoryId) where.category = { categoryId };
    if (brandId) where.brand = { brandId };
    if (search) where.name = { $like: `%${search}%` };

    const [products, total] = await this.em.findAndCount(Product, where, {
      populate: ['brand', 'category'],
      limit,
      offset: (page - 1) * limit,
      orderBy: { productId: 'DESC' },
    });

    return {
      data: products.map((p) => this.toDomain(p)),
      meta: createPaginationMeta(page, limit, total),
    };
  }

  async findOne(id: number): Promise<ProductDetailDomain> {
    const product = await this.em.findOne(
      Product,
      { productId: id },
      { populate: ['brand', 'category', 'attributes', 'marketplaceProducts.marketplace'] },
    );
    if (!product) {
      throw new NotFoundError(`Product with ID ${id} not found`);
    }
    return this.toDetailDomain(product);
  }

  async findPrices(productId: number): Promise<ProductPriceDomain[]> {
    const product = await this.em.findOne(
      Product,
      { productId },
      { populate: ['marketplaceProducts.prices', 'marketplaceProducts.marketplace'] },
    );
    if (!product) {
      throw new NotFoundError(`Product with ID ${productId} not found`);
    }
    return product.marketplaceProducts.getItems().map((mp) => {
      const prices = mp.prices
        .getItems()
        .sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime())
        .slice(0, 30)
        .map((p) => new PriceDomain(p.priceId, p.regularPrice, p.specialPrice, p.createdAt));

      return new ProductPriceDomain(
        new MarketplaceDomain(mp.marketplace.marketplaceId, mp.marketplace.name, mp.marketplace.url),
        mp.url,
        prices,
      );
    });
  }

  private toDomain(orm: Product): ProductDomain {
    return new ProductDomain(
      orm.productId,
      orm.name,
      orm.imageUrl,
      new BrandDomain(orm.brand.brandId, orm.brand.name, orm.brand.bannerUrl),
      new CategoryDomain(orm.category.categoryId, orm.category.name),
    );
  }

  private toDetailDomain(orm: Product): ProductDetailDomain {
    const attributes = orm.attributes
      .getItems()
      .map((a) => new ProductAttributeDomain(a.attributeId, a.attributeType, a.value));

    const marketplaceProducts = orm.marketplaceProducts.getItems().map(
      (mp) =>
        new MarketplaceProductDomain(
          mp.marketplaceProductId,
          new MarketplaceDomain(mp.marketplace.marketplaceId, mp.marketplace.name, mp.marketplace.url),
          mp.url,
          [],
        ),
    );

    return new ProductDetailDomain(
      orm.productId,
      orm.name,
      orm.imageUrl,
      new BrandDomain(orm.brand.brandId, orm.brand.name, orm.brand.bannerUrl),
      new CategoryDomain(orm.category.categoryId, orm.category.name),
      attributes,
      marketplaceProducts,
    );
  }
}
