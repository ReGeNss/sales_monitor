import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { MarketplaceResponseDto } from '../../marketplaces/dto/marketplace-response.dto';

export class PriceItemDto {
  @ApiProperty({ example: 1, description: 'Price record ID' })
  priceId!: number;

  @ApiProperty({ example: 1299.99, description: 'Regular price' })
  regularPrice!: number;

  @ApiPropertyOptional({ example: 999.99, description: 'Discounted price' })
  discountPrice?: number;

  @ApiProperty({ example: '2026-02-15T12:00:00.000Z', description: 'Price record creation date' })
  createdAt!: Date;
}

export class MarketplaceProductPriceDto {
  @ApiProperty({ example: 1, description: 'Marketplace product record ID' })
  marketplaceProductId!: number;

  @ApiProperty({ type: MarketplaceResponseDto, description: 'Marketplace details' })
  marketplace!: MarketplaceResponseDto;

  @ApiProperty({ example: 'https://rozetka.com.ua/product/123', description: 'Product URL on the marketplace' })
  url!: string;

  @ApiProperty({ type: [PriceItemDto], description: 'Array of price records' })
  prices!: PriceItemDto[];
}

export class PriceWithProductDto extends PriceItemDto {
  @ApiProperty({ type: MarketplaceProductPriceDto, description: 'Marketplace product details' })
  marketplaceProduct!: MarketplaceProductPriceDto;
}

export class ProductPriceHistoryItemDto {
  @ApiProperty({ type: MarketplaceResponseDto, description: 'Marketplace details' })
  marketplace!: MarketplaceResponseDto;

  @ApiProperty({ example: 'https://rozetka.com.ua/product/123', description: 'Product URL on the marketplace' })
  url!: string;

  @ApiProperty({ type: [PriceItemDto], description: 'Price history (last 30 records)' })
  prices!: PriceItemDto[];
}
