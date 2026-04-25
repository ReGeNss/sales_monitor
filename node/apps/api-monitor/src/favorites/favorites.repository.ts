import { Injectable } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { User, Product, Brand } from '@sales-monitor/database';
import { ProductDomain } from '../products/domain/product.domain';
import { BrandDomain } from '../brands/domain/brand.domain';
import { CategoryDomain } from '../categories/domain/category.domain';
import { NotFoundError } from '../common/errors';

@Injectable()
export class FavoritesRepository {
  constructor(private readonly em: EntityManager) {}

  async getFavoriteProducts(userId: number): Promise<ProductDomain[]> {
    const user = await this.em.findOne(
      User,
      { userId },
      { populate: ['favoriteProducts.brand', 'favoriteProducts.category'] },
    );
    if (!user) {
      throw new NotFoundError('User not found');
    }
    return user.favoriteProducts.getItems().map(
      (p) =>
        new ProductDomain(
          p.productId,
          p.name,
          p.imageUrl,
          new BrandDomain(p.brand.brandId, p.brand.name, p.brand.bannerUrl),
          new CategoryDomain(p.category.categoryId, p.category.name),
        ),
    );
  }

  async addFavoriteProduct(userId: number, productId: number): Promise<void> {
    const user = await this.em.findOne(User, { userId }, { populate: ['favoriteProducts'] });
    if (!user) throw new NotFoundError('User not found');

    const product = await this.em.findOne(Product, { productId });
    if (!product) throw new NotFoundError('Product not found');

    if (!user.favoriteProducts.contains(product)) {
      user.favoriteProducts.add(product);
      await this.em.flush();
    }
  }

  async removeFavoriteProduct(userId: number, productId: number): Promise<void> {
    const user = await this.em.findOne(User, { userId }, { populate: ['favoriteProducts'] });
    if (!user) throw new NotFoundError('User not found');

    const product = await this.em.findOne(Product, { productId });
    if (!product) throw new NotFoundError('Product not found');

    user.favoriteProducts.remove(product);
    await this.em.flush();
  }

  async getFavoriteBrands(userId: number): Promise<BrandDomain[]> {
    const user = await this.em.findOne(User, { userId }, { populate: ['favoriteBrands'] });
    if (!user) throw new NotFoundError('User not found');

    return user.favoriteBrands.getItems().map(
      (b) => new BrandDomain(b.brandId, b.name, b.bannerUrl),
    );
  }

  async addFavoriteBrand(userId: number, brandId: number): Promise<void> {
    const user = await this.em.findOne(User, { userId }, { populate: ['favoriteBrands'] });
    if (!user) throw new NotFoundError('User not found');

    const brand = await this.em.findOne(Brand, { brandId });
    if (!brand) throw new NotFoundError('Brand not found');

    if (!user.favoriteBrands.contains(brand)) {
      user.favoriteBrands.add(brand);
      await this.em.flush();
    }
  }

  async removeFavoriteBrand(userId: number, brandId: number): Promise<void> {
    const user = await this.em.findOne(User, { userId }, { populate: ['favoriteBrands'] });
    if (!user) throw new NotFoundError('User not found');

    const brand = await this.em.findOne(Brand, { brandId });
    if (!brand) throw new NotFoundError('Brand not found');

    user.favoriteBrands.remove(brand);
    await this.em.flush();
  }
}
