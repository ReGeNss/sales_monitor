import { Injectable, NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { User, Product, Brand } from '@sales-monitor/database';

@Injectable()
export class FavoritesService {
  constructor(private readonly em: EntityManager) {}

  async getFavoriteProducts(userId: number) {
    const user = await this.em.findOne(
      User,
      { userId },
      { populate: ['favoriteProducts.brand', 'favoriteProducts.category'] },
    );

    if (!user) {
      throw new NotFoundException('User not found');
    }

    return user.favoriteProducts.getItems();
  }

  async addFavoriteProduct(userId: number, productId: number) {
    const user = await this.em.findOne(User, { userId }, { populate: ['favoriteProducts'] });
    if (!user) {
      throw new NotFoundException('User not found');
    }

    const product = await this.em.findOne(Product, { productId });
    if (!product) {
      throw new NotFoundException('Product not found');
    }

    if (!user.favoriteProducts.contains(product)) {
      user.favoriteProducts.add(product);
      await this.em.flush();
    }

    return { message: 'Product added to favorites' };
  }

  async removeFavoriteProduct(userId: number, productId: number) {
    const user = await this.em.findOne(User, { userId }, { populate: ['favoriteProducts'] });
    if (!user) {
      throw new NotFoundException('User not found');
    }

    const product = await this.em.findOne(Product, { productId });
    if (!product) {
      throw new NotFoundException('Product not found');
    }

    user.favoriteProducts.remove(product);
    await this.em.flush();

    return { message: 'Product removed from favorites' };
  }

  async getFavoriteBrands(userId: number) {
    const user = await this.em.findOne(User, { userId }, { populate: ['favoriteBrands'] });

    if (!user) {
      throw new NotFoundException('User not found');
    }

    return user.favoriteBrands.getItems();
  }

  async addFavoriteBrand(userId: number, brandId: number) {
    const user = await this.em.findOne(User, { userId }, { populate: ['favoriteBrands'] });
    if (!user) {
      throw new NotFoundException('User not found');
    }

    const brand = await this.em.findOne(Brand, { brandId });
    if (!brand) {
      throw new NotFoundException('Brand not found');
    }

    if (!user.favoriteBrands.contains(brand)) {
      user.favoriteBrands.add(brand);
      await this.em.flush();
    }

    return { message: 'Brand added to favorites' };
  }

  async removeFavoriteBrand(userId: number, brandId: number) {
    const user = await this.em.findOne(User, { userId }, { populate: ['favoriteBrands'] });
    if (!user) {
      throw new NotFoundException('User not found');
    }

    const brand = await this.em.findOne(Brand, { brandId });
    if (!brand) {
      throw new NotFoundException('Brand not found');
    }

    user.favoriteBrands.remove(brand);
    await this.em.flush();

    return { message: 'Brand removed from favorites' };
  }
}
