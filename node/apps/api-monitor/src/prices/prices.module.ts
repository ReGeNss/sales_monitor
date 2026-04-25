import { Module } from '@nestjs/common';
import { MikroOrmModule } from '@mikro-orm/nestjs';
import { Price } from '@sales-monitor/database';
import { PricesService } from './prices.service';
import { PricesController } from './prices.controller';
import { PricesRepository } from './prices.repository';

@Module({
  imports: [MikroOrmModule.forFeature([Price])],
  providers: [PricesService, PricesRepository],
  controllers: [PricesController],
})
export class PricesModule {}
