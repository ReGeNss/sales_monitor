import { ApiProperty } from '@nestjs/swagger';

export class CategoryResponseDto {
  @ApiProperty({ example: 1, description: 'Category ID' })
  categoryId!: number;

  @ApiProperty({ example: 'Electronics', description: 'Category name' })
  name!: string;
}
