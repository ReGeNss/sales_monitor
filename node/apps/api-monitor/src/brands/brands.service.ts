import { Injectable } from '@nestjs/common';
import { createPaginationMeta } from '@sales-monitor/common';
import { BrandsRepository } from './brands.repository';
import { BrandDomain } from './domain/brand.domain';
import { ProductDomain } from '../products/domain/product.domain';

@Injectable()
export class BrandsService {
  constructor(private readonly brandsRepository: BrandsRepository) {}

  async findAll(): Promise<BrandDomain[]> {
    return this.brandsRepository.findAll();
  }

  async findOne(id: number): Promise<BrandDomain> {
    return this.brandsRepository.findOne(id);
  }

  async getBrandProducts(brandId: number, page = 1, limit = 20) {
    const products = await this.brandsRepository.findProducts(brandId);
    const total = products.length;
    const data: ProductDomain[] = products.slice((page - 1) * limit, page * limit);

    return {
      data,
      meta: createPaginationMeta(page, limit, total),
    };
  }
}
