import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class BrandResponseDto {
  @ApiProperty({ example: 1, description: 'Brand ID' })
  brandId!: number;

  @ApiProperty({ example: 'Samsung', description: 'Brand name' })
  name!: string;

  @ApiPropertyOptional({ example: 'https://example.com/banner.jpg', description: 'Brand banner image URL' })
  bannerUrl?: string;
}
