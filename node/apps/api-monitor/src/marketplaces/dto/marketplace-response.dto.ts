import { ApiProperty } from '@nestjs/swagger';

export class MarketplaceResponseDto {
  @ApiProperty({ example: 1, description: 'Marketplace ID' })
  marketplaceId!: number;

  @ApiProperty({ example: 'Rozetka', description: 'Marketplace name' })
  name!: string;

  @ApiProperty({ example: 'https://rozetka.com.ua', description: 'Marketplace website URL' })
  url!: string;
}
