import { Injectable } from '@nestjs/common';
import { FavoritesRepository } from './favorites.repository';
import { ProductDomain } from '../products/domain/product.domain';
import { BrandDomain } from '../brands/domain/brand.domain';

@Injectable()
export class FavoritesService {
  constructor(private readonly favoritesRepository: FavoritesRepository) {}

  async getFavoriteProducts(userId: number): Promise<ProductDomain[]> {
    return this.favoritesRepository.getFavoriteProducts(userId);
  }

  async addFavoriteProduct(userId: number, productId: number) {
    await this.favoritesRepository.addFavoriteProduct(userId, productId);
    return { message: 'Product added to favorites' };
  }

  async removeFavoriteProduct(userId: number, productId: number) {
    await this.favoritesRepository.removeFavoriteProduct(userId, productId);
    return { message: 'Product removed from favorites' };
  }

  async getFavoriteBrands(userId: number): Promise<BrandDomain[]> {
    return this.favoritesRepository.getFavoriteBrands(userId);
  }

  async addFavoriteBrand(userId: number, brandId: number) {
    await this.favoritesRepository.addFavoriteBrand(userId, brandId);
    return { message: 'Brand added to favorites' };
  }

  async removeFavoriteBrand(userId: number, brandId: number) {
    await this.favoritesRepository.removeFavoriteBrand(userId, brandId);
    return { message: 'Brand removed from favorites' };
  }
}
