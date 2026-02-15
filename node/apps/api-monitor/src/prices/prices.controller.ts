import { Controller, Get, Query, ParseIntPipe } from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiQuery,
  ApiBearerAuth,
  ApiUnauthorizedResponse,
} from '@nestjs/swagger';
import { PricesService } from './prices.service';
import { PriceItemDto } from './dto/price-response.dto';

@ApiTags('prices')
@Controller('prices')
@ApiBearerAuth()
@ApiUnauthorizedResponse({ description: 'Authorization required (Bearer token)' })
export class PricesController {
  constructor(private readonly pricesService: PricesService) {}

  @Get('latest')
  @ApiOperation({
    summary: 'Get latest prices',
    description: 'Returns the most recent price records with product and marketplace information, sorted by date (newest first).',
  })
  @ApiQuery({
    name: 'limit',
    required: false,
    type: Number,
    description: 'Maximum number of records (default: 100)',
    example: 50,
  })
  @ApiResponse({
    status: 200,
    description: 'Array of latest price records',
    type: [PriceItemDto],
  })
  async getLatestPrices(
    @Query('limit', new ParseIntPipe({ optional: true })) limit?: number,
  ) {
    return this.pricesService.getLatestPrices(limit);
  }

  @Get('trends')
  @ApiOperation({
    summary: 'Get price trends',
    description: 'Returns price history for a specified period (default: 30 days). Can be filtered by a specific product.',
  })
  @ApiQuery({
    name: 'productId',
    required: false,
    type: Number,
    description: 'Filter by product ID',
    example: 1,
  })
  @ApiQuery({
    name: 'days',
    required: false,
    type: Number,
    description: 'Number of days to analyze (default: 30)',
    example: 30,
  })
  @ApiResponse({
    status: 200,
    description: 'Array of price records for the specified period',
    type: [PriceItemDto],
  })
  async getPriceTrends(
    @Query('productId', new ParseIntPipe({ optional: true })) productId?: number,
    @Query('days', new ParseIntPipe({ optional: true })) days?: number,
  ) {
    return this.pricesService.getPriceTrends(productId, days);
  }
}
