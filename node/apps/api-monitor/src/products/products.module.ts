import { Module } from '@nestjs/common';
import { MikroOrmModule } from '@mikro-orm/nestjs';
import { Product } from '@sales-monitor/database';
import { ProductsService } from './products.service';
import { ProductsController } from './products.controller';

@Module({
  imports: [MikroOrmModule.forFeature([Product])],
  providers: [ProductsService],
  controllers: [ProductsController],
})
export class ProductsModule {}
