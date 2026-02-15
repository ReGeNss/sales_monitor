import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { BrandResponseDto } from '../../brands/dto/brand-response.dto';
import { CategoryResponseDto } from '../../categories/dto/category-response.dto';
import { PaginationMetaDto } from '../../common/dto/pagination-meta.dto';

export class ProductAttributeResponseDto {
  @ApiProperty({ example: 1, description: 'Attribute ID' })
  attributeId!: number;

  @ApiProperty({ example: 'volume', description: "Attribute type ('volume' or 'weight')" })
  attributeType!: string;

  @ApiProperty({ example: '500ml', description: 'Attribute value' })
  value!: string;
}

export class MarketplaceProductResponseDto {
  @ApiProperty({ example: 1, description: 'Marketplace product record ID' })
  marketplaceProductId!: number;

  @ApiProperty({ example: 'https://rozetka.com.ua/product/123', description: 'Product URL on the marketplace' })
  url!: string;
}

export class ProductResponseDto {
  @ApiProperty({ example: 1, description: 'Product ID' })
  productId!: number;

  @ApiProperty({ example: 'Head & Shoulders Shampoo 400ml', description: 'Product name' })
  name!: string;

  @ApiPropertyOptional({ example: 'https://example.com/image.jpg', description: 'Product image URL' })
  imageUrl?: string;

  @ApiProperty({ type: BrandResponseDto, description: 'Product brand' })
  brand!: BrandResponseDto;

  @ApiProperty({ type: CategoryResponseDto, description: 'Product category' })
  category!: CategoryResponseDto;
}

export class ProductDetailResponseDto extends ProductResponseDto {
  @ApiProperty({ type: [ProductAttributeResponseDto], description: 'Product attributes' })
  attributes!: ProductAttributeResponseDto[];

  @ApiProperty({ type: [MarketplaceProductResponseDto], description: 'Product listings on marketplaces' })
  marketplaceProducts!: MarketplaceProductResponseDto[];
}

export class PaginatedProductResponseDto {
  @ApiProperty({ type: [ProductResponseDto], description: 'Array of products' })
  data!: ProductResponseDto[];

  @ApiProperty({ type: PaginationMetaDto, description: 'Pagination metadata' })
  meta!: PaginationMetaDto;
}
