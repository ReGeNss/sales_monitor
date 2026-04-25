import { Module } from '@nestjs/common';
import { MikroOrmModule } from '@mikro-orm/nestjs';
import { Marketplace } from '@sales-monitor/database';
import { MarketplacesService } from './marketplaces.service';
import { MarketplacesController } from './marketplaces.controller';
import { MarketplacesRepository } from './marketplaces.repository';

@Module({
  imports: [MikroOrmModule.forFeature([Marketplace])],
  providers: [MarketplacesService, MarketplacesRepository],
  controllers: [MarketplacesController],
})
export class MarketplacesModule {}
