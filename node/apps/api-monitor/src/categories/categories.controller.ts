import { Controller, Get, Param, Query, ParseIntPipe } from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiParam,
  ApiQuery,
  ApiBearerAuth,
  ApiUnauthorizedResponse,
  ApiNotFoundResponse,
} from '@nestjs/swagger';
import { CategoriesService } from './categories.service';
import { CategoryResponseDto } from './dto/category-response.dto';
import { PaginatedCategoryProductsResponseDto } from './dto/paginated-category-products-response.dto';

@ApiTags('categories')
@Controller('categories')
@ApiBearerAuth()
@ApiUnauthorizedResponse({ description: 'Authorization required (Bearer token)' })
export class CategoriesController {
  constructor(private readonly categoriesService: CategoriesService) {}

  @Get()
  @ApiOperation({
    summary: 'Get all categories',
    description: 'Returns a list of all product categories sorted by name.',
  })
  @ApiResponse({
    status: 200,
    description: 'List of categories',
    type: [CategoryResponseDto],
  })
  async findAll() {
    return this.categoriesService.findAll();
  }

  @Get(':id')
  @ApiOperation({
    summary: 'Get category by ID',
    description: 'Returns information about a specific category.',
  })
  @ApiParam({ name: 'id', type: Number, description: 'Category ID', example: 1 })
  @ApiResponse({
    status: 200,
    description: 'Category details',
    type: CategoryResponseDto,
  })
  @ApiNotFoundResponse({ description: 'Category not found' })
  async findOne(@Param('id', ParseIntPipe) id: number) {
    return this.categoriesService.findOne(id);
  }

  @Get(':id/products')
  @ApiOperation({
    summary: 'Get products in category',
    description: 'Returns a paginated list of products belonging to the specified category.',
  })
  @ApiParam({ name: 'id', type: Number, description: 'Category ID', example: 1 })
  @ApiQuery({ name: 'page', required: false, type: Number, description: 'Page number (default: 1)', example: 1 })
  @ApiQuery({ name: 'limit', required: false, type: Number, description: 'Items per page (default: 20)', example: 20 })
  @ApiResponse({
    status: 200,
    description: 'Paginated list of products in the category',
    type: PaginatedCategoryProductsResponseDto,
  })
  @ApiNotFoundResponse({ description: 'Category not found' })
  async getCategoryProducts(
    @Param('id', ParseIntPipe) id: number,
    @Query('page', new ParseIntPipe({ optional: true })) page?: number,
    @Query('limit', new ParseIntPipe({ optional: true })) limit?: number,
  ) {
    return this.categoriesService.getCategoryProducts(id, page, limit);
  }
}
