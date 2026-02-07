import {
  Collection,
  Entity,
  ManyToMany,
  PrimaryKey,
  Property,
} from "@mikro-orm/core";
import { Brand } from "./brand.entity";
import { Product } from "./product.entity";

@Entity({ tableName: "users" })
export class User {
  @PrimaryKey({ fieldName: "user_id", autoincrement: true })
  userId!: number;

  @Property({ fieldName: "login", length: 255, unique: true })
  login!: string;

  @Property({ fieldName: "password", length: 255, hidden: true })
  password!: string;

  @Property({ fieldName: "nf_token", type: "text", nullable: true })
  nfToken?: string;

  @ManyToMany(() => Product, undefined, {
    pivotTable: "favorite_products",
    joinColumn: "user_id",
    inverseJoinColumn: "product_id",
  })
  favoriteProducts = new Collection<Product>(this);

  @ManyToMany(() => Brand, undefined, {
    pivotTable: "favorite_brands",
    joinColumn: "user_id",
    inverseJoinColumn: "brand_id",
  })
  favoriteBrands = new Collection<Brand>(this);
}
