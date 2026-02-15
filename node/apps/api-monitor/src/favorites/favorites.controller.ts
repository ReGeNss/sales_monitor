import { Controller, Get, Post, Delete, Param, ParseIntPipe, UseGuards } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth, ApiParam } from '@nestjs/swagger';
import { FavoritesService } from './favorites.service';
import { RegisteredUserGuard } from '../auth/guards/registered-user.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';
import { User } from '@sales-monitor/database';

@ApiTags('favorites')
@Controller('favorites')
@UseGuards(RegisteredUserGuard)
@ApiBearerAuth()
export class FavoritesController {
  constructor(private readonly favoritesService: FavoritesService) {}

  @Get('products')
  @ApiOperation({ summary: 'Get user favorite products' })
  @ApiResponse({ status: 200, description: 'List of favorite products' })
  async getFavoriteProducts(@CurrentUser() user: User) {
    return this.favoritesService.getFavoriteProducts(user.userId);
  }

  @Post('products/:productId')
  @ApiOperation({ summary: 'Add product to favorites' })
  @ApiParam({ name: 'productId', type: 'number' })
  @ApiResponse({ status: 201, description: 'Product added to favorites' })
  async addFavoriteProduct(
    @CurrentUser() user: User,
    @Param('productId', ParseIntPipe) productId: number,
  ) {
    return this.favoritesService.addFavoriteProduct(user.userId, productId);
  }

  @Delete('products/:productId')
  @ApiOperation({ summary: 'Remove product from favorites' })
  @ApiParam({ name: 'productId', type: 'number' })
  @ApiResponse({ status: 200, description: 'Product removed from favorites' })
  async removeFavoriteProduct(
    @CurrentUser() user: User,
    @Param('productId', ParseIntPipe) productId: number,
  ) {
    return this.favoritesService.removeFavoriteProduct(user.userId, productId);
  }

  @Get('brands')
  @ApiOperation({ summary: 'Get user favorite brands' })
  @ApiResponse({ status: 200, description: 'List of favorite brands' })
  async getFavoriteBrands(@CurrentUser() user: User) {
    return this.favoritesService.getFavoriteBrands(user.userId);
  }

  @Post('brands/:brandId')
  @ApiOperation({ summary: 'Add brand to favorites' })
  @ApiParam({ name: 'brandId', type: 'number' })
  @ApiResponse({ status: 201, description: 'Brand added to favorites' })
  async addFavoriteBrand(
    @CurrentUser() user: User,
    @Param('brandId', ParseIntPipe) brandId: number,
  ) {
    return this.favoritesService.addFavoriteBrand(user.userId, brandId);
  }

  @Delete('brands/:brandId')
  @ApiOperation({ summary: 'Remove brand from favorites' })
  @ApiParam({ name: 'brandId', type: 'number' })
  @ApiResponse({ status: 200, description: 'Brand removed from favorites' })
  async removeFavoriteBrand(
    @CurrentUser() user: User,
    @Param('brandId', ParseIntPipe) brandId: number,
  ) {
    return this.favoritesService.removeFavoriteBrand(user.userId, brandId);
  }
}
