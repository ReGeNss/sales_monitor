import { Injectable } from '@nestjs/common';
import { PaginatedResponse } from '@sales-monitor/common';
import { ProductsRepository } from './products.repository';
import { ProductFilterDto } from './dto/product-filter.dto';
import { ProductDomain } from './domain/product.domain';
import { ProductDetailDomain } from './domain/product-detail.domain';
import { ProductPriceDomain } from './domain/product-price.domain';

@Injectable()
export class ProductsService {
  constructor(private readonly productsRepository: ProductsRepository) {}

  async findAll(filterDto: ProductFilterDto): Promise<PaginatedResponse<ProductDomain>> {
    return this.productsRepository.findAll(filterDto);
  }

  async findOne(id: number): Promise<ProductDetailDomain> {
    return this.productsRepository.findOne(id);
  }

  async getProductPrices(productId: number): Promise<ProductPriceDomain[]> {
    return this.productsRepository.findPrices(productId);
  }
}
