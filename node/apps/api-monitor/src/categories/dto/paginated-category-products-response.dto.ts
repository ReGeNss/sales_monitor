import { ApiProperty } from '@nestjs/swagger';
import { ProductResponseDto } from '../../products/dto/product-response.dto';
import { PaginationMetaDto } from '../../common/dto/pagination-meta.dto';

export class PaginatedCategoryProductsResponseDto {
  @ApiProperty({ type: [ProductResponseDto], description: 'Array of products in the category' })
  data!: ProductResponseDto[];

  @ApiProperty({ type: PaginationMetaDto, description: 'Pagination metadata' })
  meta!: PaginationMetaDto;
}
