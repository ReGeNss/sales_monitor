import { ApiProperty } from '@nestjs/swagger';

export class PaginationMetaDto {
  @ApiProperty({ example: 1, description: 'Current page number' })
  page!: number;

  @ApiProperty({ example: 20, description: 'Number of items per page' })
  limit!: number;

  @ApiProperty({ example: 150, description: 'Total number of items' })
  total!: number;

  @ApiProperty({ example: 8, description: 'Total number of pages' })
  totalPages!: number;

  @ApiProperty({ example: true, description: 'Whether a next page exists' })
  hasNext!: boolean;

  @ApiProperty({ example: false, description: 'Whether a previous page exists' })
  hasPrev!: boolean;
}
