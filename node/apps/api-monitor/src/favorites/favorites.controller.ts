import { Controller, Get, Post, Delete, Param, ParseIntPipe, UseGuards } from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
  ApiParam,
  ApiUnauthorizedResponse,
  ApiForbiddenResponse,
  ApiNotFoundResponse,
} from '@nestjs/swagger';
import { FavoritesService } from './favorites.service';
import { RegisteredUserGuard } from '../auth/guards/registered-user.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';
import { User } from '@sales-monitor/database';
import { ProductResponseDto } from '../products/dto/product-response.dto';
import { BrandResponseDto } from '../brands/dto/brand-response.dto';
import { MessageResponseDto } from '../common/dto/message-response.dto';

@ApiTags('favorites')
@Controller('favorites')
@UseGuards(RegisteredUserGuard)
@ApiBearerAuth()
@ApiUnauthorizedResponse({ description: 'Authorization required (Bearer token)' })
@ApiForbiddenResponse({ description: 'Access restricted to registered users only (not guests)' })
export class FavoritesController {
  constructor(private readonly favoritesService: FavoritesService) {}

  // ── Products ──────────────────────────────────────────────

  @Get('products')
  @ApiOperation({
    summary: 'Get favorite products',
    description: 'Returns a list of products added to favorites by the current user.',
  })
  @ApiResponse({
    status: 200,
    description: 'List of favorite products',
    type: [ProductResponseDto],
  })
  async getFavoriteProducts(@CurrentUser() user: User) {
    return this.favoritesService.getFavoriteProducts(user.userId);
  }

  @Post('products/:productId')
  @ApiOperation({
    summary: 'Add product to favorites',
    description: 'Adds the specified product to the current user\'s favorites list. If the product is already in favorites, no duplicate is created.',
  })
  @ApiParam({ name: 'productId', type: Number, description: 'Product ID', example: 1 })
  @ApiResponse({
    status: 201,
    description: 'Product added to favorites',
    type: MessageResponseDto,
  })
  @ApiNotFoundResponse({ description: 'Product or user not found' })
  async addFavoriteProduct(
    @CurrentUser() user: User,
    @Param('productId', ParseIntPipe) productId: number,
  ) {
    return this.favoritesService.addFavoriteProduct(user.userId, productId);
  }

  @Delete('products/:productId')
  @ApiOperation({
    summary: 'Remove product from favorites',
    description: 'Removes the specified product from the current user\'s favorites list.',
  })
  @ApiParam({ name: 'productId', type: Number, description: 'Product ID', example: 1 })
  @ApiResponse({
    status: 200,
    description: 'Product removed from favorites',
    type: MessageResponseDto,
  })
  @ApiNotFoundResponse({ description: 'Product or user not found' })
  async removeFavoriteProduct(
    @CurrentUser() user: User,
    @Param('productId', ParseIntPipe) productId: number,
  ) {
    return this.favoritesService.removeFavoriteProduct(user.userId, productId);
  }

  // ── Brands ────────────────────────────────────────────────

  @Get('brands')
  @ApiOperation({
    summary: 'Get favorite brands',
    description: 'Returns a list of brands added to favorites by the current user.',
  })
  @ApiResponse({
    status: 200,
    description: 'List of favorite brands',
    type: [BrandResponseDto],
  })
  async getFavoriteBrands(@CurrentUser() user: User) {
    return this.favoritesService.getFavoriteBrands(user.userId);
  }

  @Post('brands/:brandId')
  @ApiOperation({
    summary: 'Add brand to favorites',
    description: 'Adds the specified brand to the current user\'s favorites list. If the brand is already in favorites, no duplicate is created.',
  })
  @ApiParam({ name: 'brandId', type: Number, description: 'Brand ID', example: 1 })
  @ApiResponse({
    status: 201,
    description: 'Brand added to favorites',
    type: MessageResponseDto,
  })
  @ApiNotFoundResponse({ description: 'Brand or user not found' })
  async addFavoriteBrand(
    @CurrentUser() user: User,
    @Param('brandId', ParseIntPipe) brandId: number,
  ) {
    return this.favoritesService.addFavoriteBrand(user.userId, brandId);
  }

  @Delete('brands/:brandId')
  @ApiOperation({
    summary: 'Remove brand from favorites',
    description: 'Removes the specified brand from the current user\'s favorites list.',
  })
  @ApiParam({ name: 'brandId', type: Number, description: 'Brand ID', example: 1 })
  @ApiResponse({
    status: 200,
    description: 'Brand removed from favorites',
    type: MessageResponseDto,
  })
  @ApiNotFoundResponse({ description: 'Brand or user not found' })
  async removeFavoriteBrand(
    @CurrentUser() user: User,
    @Param('brandId', ParseIntPipe) brandId: number,
  ) {
    return this.favoritesService.removeFavoriteBrand(user.userId, brandId);
  }
}
