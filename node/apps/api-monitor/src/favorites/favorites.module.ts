import { Module } from '@nestjs/common';
import { MikroOrmModule } from '@mikro-orm/nestjs';
import { User, Product, Brand } from '@sales-monitor/database';
import { FavoritesService } from './favorites.service';
import { FavoritesController } from './favorites.controller';
import { FavoritesRepository } from './favorites.repository';

@Module({
  imports: [MikroOrmModule.forFeature([User, Product, Brand])],
  providers: [FavoritesService, FavoritesRepository],
  controllers: [FavoritesController],
})
export class FavoritesModule {}
