import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { MikroOrmModule } from "@mikro-orm/nestjs";
import { MySqlDriver } from "@mikro-orm/mysql";
import * as path from "path";
import * as fs from "fs";
import { AuthModule } from "./auth/auth.module";
import { UsersModule } from "./users/users.module";
import { FavoritesModule } from "./favorites/favorites.module";
import { ProductsModule } from "./products/products.module";
import { PricesModule } from "./prices/prices.module";
import { CategoriesModule } from "./categories/categories.module";
import { BrandsModule } from "./brands/brands.module";
import { MarketplacesModule } from "./marketplaces/marketplaces.module";

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: path.resolve(process.cwd(), ".env"),
    }),
    MikroOrmModule.forRootAsync({
      useFactory: async () => {
        const {
          User,
          Product,
          Price,
          Category,
          Brand,
          Marketplace,
          MarketplaceProduct,
          ProductAttribute,
        } = await import("@sales-monitor/database");
        const isDocker = fs.existsSync("/.dockerenv");
        return {
          driver: MySqlDriver,
          entities: [
            User,
            Product,
            Price,
            Category,
            Brand,
            Marketplace,
            MarketplaceProduct,
            ProductAttribute,
          ],
          host: isDocker ? process.env.DB_HOST : "localhost",
          port: parseInt(process.env.DB_PORT as string),
          dbName: process.env.DATABASE_NAME,
          user: process.env.DB_USER_NAME,
          password: process.env.DB_USER_PASSWORD,
          debug: process.env.NODE_ENV === "development",
          allowGlobalContext: true,
        };
      },
    }),
    AuthModule,
    UsersModule,
    FavoritesModule,
    ProductsModule,
    PricesModule,
    CategoriesModule,
    BrandsModule,
    MarketplacesModule,
  ],
})
export class AppModule {}
