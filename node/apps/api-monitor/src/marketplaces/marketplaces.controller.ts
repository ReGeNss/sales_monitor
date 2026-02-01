import { Controller, Get, Param, ParseIntPipe } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiParam } from '@nestjs/swagger';
import { MarketplacesService } from './marketplaces.service';

@ApiTags('marketplaces')
@Controller('marketplaces')
export class MarketplacesController {
  constructor(private readonly marketplacesService: MarketplacesService) {}

  @Get()
  @ApiOperation({ summary: 'Get all marketplaces' })
  @ApiResponse({ status: 200, description: 'List of marketplaces' })
  async findAll() {
    return this.marketplacesService.findAll();
  }

  @Get(':id')
  @ApiOperation({ summary: 'Get marketplace by ID' })
  @ApiParam({ name: 'id', type: 'number' })
  @ApiResponse({ status: 200, description: 'Marketplace details' })
  @ApiResponse({ status: 404, description: 'Marketplace not found' })
  async findOne(@Param('id', ParseIntPipe) id: number) {
    return this.marketplacesService.findOne(id);
  }
}
