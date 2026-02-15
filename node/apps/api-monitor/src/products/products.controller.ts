import { Controller, Get, Param, Query, ParseIntPipe } from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiParam,
  ApiBearerAuth,
  ApiUnauthorizedResponse,
  ApiNotFoundResponse,
} from '@nestjs/swagger';
import { ProductsService } from './products.service';
import { ProductFilterDto } from './dto/product-filter.dto';
import {
  PaginatedProductResponseDto,
  ProductDetailResponseDto,
} from './dto/product-response.dto';
import { ProductPriceHistoryItemDto } from '../prices/dto/price-response.dto';

@ApiTags('products')
@Controller('products')
@ApiBearerAuth()
@ApiUnauthorizedResponse({ description: 'Authorization required (Bearer token)' })
export class ProductsController {
  constructor(private readonly productsService: ProductsService) {}

  @Get()
  @ApiOperation({
    summary: 'Get all products with filtering and pagination',
    description: 'Returns a paginated list of products with optional filtering by category, brand, and search by name.',
  })
  @ApiResponse({
    status: 200,
    description: 'Paginated list of products',
    type: PaginatedProductResponseDto,
  })
  async findAll(@Query() filterDto: ProductFilterDto) {
    return this.productsService.findAll(filterDto);
  }

  @Get(':id')
  @ApiOperation({
    summary: 'Get product by ID',
    description: 'Returns detailed product information including brand, category, attributes, and marketplace listings.',
  })
  @ApiParam({ name: 'id', type: Number, description: 'Product ID', example: 1 })
  @ApiResponse({
    status: 200,
    description: 'Product details',
    type: ProductDetailResponseDto,
  })
  @ApiNotFoundResponse({ description: 'Product not found' })
  async findOne(@Param('id', ParseIntPipe) id: number) {
    return this.productsService.findOne(id);
  }

  @Get(':id/prices')
  @ApiOperation({
    summary: 'Get product price history',
    description: 'Returns the price history of a product across marketplaces (last 30 records per marketplace).',
  })
  @ApiParam({ name: 'id', type: Number, description: 'Product ID', example: 1 })
  @ApiResponse({
    status: 200,
    description: 'Product price history grouped by marketplace',
    type: [ProductPriceHistoryItemDto],
  })
  @ApiNotFoundResponse({ description: 'Product not found' })
  async getProductPrices(@Param('id', ParseIntPipe) id: number) {
    return this.productsService.getProductPrices(id);
  }
}
