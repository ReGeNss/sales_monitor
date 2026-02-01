import { Controller, Get, Query, ParseIntPipe } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiQuery } from '@nestjs/swagger';
import { PricesService } from './prices.service';

@ApiTags('prices')
@Controller('prices')
export class PricesController {
  constructor(private readonly pricesService: PricesService) {}

  @Get('latest')
  @ApiOperation({ summary: 'Get latest prices' })
  @ApiQuery({ name: 'limit', required: false, type: 'number', description: 'Number of results' })
  @ApiResponse({ status: 200, description: 'Latest prices' })
  async getLatestPrices(
    @Query('limit', new ParseIntPipe({ optional: true })) limit?: number,
  ) {
    return this.pricesService.getLatestPrices(limit);
  }

  @Get('trends')
  @ApiOperation({ summary: 'Get price trends' })
  @ApiQuery({ name: 'productId', required: false, type: 'number' })
  @ApiQuery({ name: 'days', required: false, type: 'number', description: 'Number of days to analyze' })
  @ApiResponse({ status: 200, description: 'Price trends' })
  async getPriceTrends(
    @Query('productId', new ParseIntPipe({ optional: true })) productId?: number,
    @Query('days', new ParseIntPipe({ optional: true })) days?: number,
  ) {
    return this.pricesService.getPriceTrends(productId, days);
  }
}
