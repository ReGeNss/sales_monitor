import { Module } from '@nestjs/common';
import { MikroOrmModule } from '@mikro-orm/nestjs';
import { Brand } from '@sales-monitor/database';
import { BrandsService } from './brands.service';
import { BrandsController } from './brands.controller';
import { BrandsRepository } from './brands.repository';

@Module({
  imports: [MikroOrmModule.forFeature([Brand])],
  providers: [BrandsService, BrandsRepository],
  controllers: [BrandsController],
})
export class BrandsModule {}
