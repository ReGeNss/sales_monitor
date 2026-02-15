import { Controller, Get, Param, ParseIntPipe } from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiParam,
  ApiBearerAuth,
  ApiUnauthorizedResponse,
  ApiNotFoundResponse,
} from '@nestjs/swagger';
import { MarketplacesService } from './marketplaces.service';
import { MarketplaceResponseDto } from './dto/marketplace-response.dto';

@ApiTags('marketplaces')
@Controller('marketplaces')
@ApiBearerAuth()
@ApiUnauthorizedResponse({ description: 'Authorization required (Bearer token)' })
export class MarketplacesController {
  constructor(private readonly marketplacesService: MarketplacesService) {}

  @Get()
  @ApiOperation({
    summary: 'Get all marketplaces',
    description: 'Returns a list of all marketplaces sorted by name.',
  })
  @ApiResponse({
    status: 200,
    description: 'List of marketplaces',
    type: [MarketplaceResponseDto],
  })
  async findAll() {
    return this.marketplacesService.findAll();
  }

  @Get(':id')
  @ApiOperation({
    summary: 'Get marketplace by ID',
    description: 'Returns information about a specific marketplace.',
  })
  @ApiParam({ name: 'id', type: Number, description: 'Marketplace ID', example: 1 })
  @ApiResponse({
    status: 200,
    description: 'Marketplace details',
    type: MarketplaceResponseDto,
  })
  @ApiNotFoundResponse({ description: 'Marketplace not found' })
  async findOne(@Param('id', ParseIntPipe) id: number) {
    return this.marketplacesService.findOne(id);
  }
}
