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
import { BrandsService } from './brands.service';
import { BrandResponseDto } from './dto/brand-response.dto';
import { PaginatedBrandProductsResponseDto } from './dto/paginated-brand-products-response.dto';

@ApiTags('brands')
@Controller('brands')
@ApiBearerAuth()
@ApiUnauthorizedResponse({ description: 'Authorization required (Bearer token)' })
export class BrandsController {
  constructor(private readonly brandsService: BrandsService) {}

  @Get()
  @ApiOperation({
    summary: 'Get all brands',
    description: 'Returns a list of all brands sorted by name.',
  })
  @ApiResponse({
    status: 200,
    description: 'List of brands',
    type: [BrandResponseDto],
  })
  async findAll() {
    return this.brandsService.findAll();
  }

  @Get(':id')
  @ApiOperation({
    summary: 'Get brand by ID',
    description: 'Returns information about a specific brand.',
  })
  @ApiParam({ name: 'id', type: Number, description: 'Brand ID', example: 1 })
  @ApiResponse({
    status: 200,
    description: 'Brand details',
    type: BrandResponseDto,
  })
  @ApiNotFoundResponse({ description: 'Brand not found' })
  async findOne(@Param('id', ParseIntPipe) id: number) {
    return this.brandsService.findOne(id);
  }

  @Get(':id/products')
  @ApiOperation({
    summary: 'Get products by brand',
    description: 'Returns a paginated list of products for the specified brand.',
  })
  @ApiParam({ name: 'id', type: Number, description: 'Brand ID', example: 1 })
  @ApiQuery({ name: 'page', required: false, type: Number, description: 'Page number (default: 1)', example: 1 })
  @ApiQuery({ name: 'limit', required: false, type: Number, description: 'Items per page (default: 20)', example: 20 })
  @ApiResponse({
    status: 200,
    description: 'Paginated list of products by the brand',
    type: PaginatedBrandProductsResponseDto,
  })
  @ApiNotFoundResponse({ description: 'Brand not found' })
  async getBrandProducts(
    @Param('id', ParseIntPipe) id: number,
    @Query('page', new ParseIntPipe({ optional: true })) page?: number,
    @Query('limit', new ParseIntPipe({ optional: true })) limit?: number,
  ) {
    return this.brandsService.getBrandProducts(id, page, limit);
  }
}
