import { Module } from '@nestjs/common';
import { MikroOrmModule } from '@mikro-orm/nestjs';
import { Price } from '@sales-monitor/database';
import { PricesService } from './prices.service';
import { PricesController } from './prices.controller';

@Module({
  imports: [MikroOrmModule.forFeature([Price])],
  providers: [PricesService],
  controllers: [PricesController],
})
export class PricesModule {}
