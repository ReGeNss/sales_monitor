import { Entity, PrimaryKey, Property, ManyToOne, OneToMany, ManyToMany, Collection } from '@mikro-orm/core';
import { Brand } from './brand.entity';
import { Category } from './category.entity';
import { MarketplaceProduct } from './marketplace-product.entity';
import { ProductAttribute } from './product-attribute.entity';

@Entity({ tableName: 'products' })
export class Product {
  @PrimaryKey({ fieldName: 'product_id' })
  productId!: number;

  @Property({ fieldName: 'name_fingerprint', length: 255, nullable: true, unique: true })
  nameFingerprint?: string;

  @ManyToOne(() => Brand, { fieldName: 'brand_id' })
  brand!: Brand;

  @Property({ fieldName: 'name', length: 255 })
  name!: string;

  @ManyToOne(() => Category, { fieldName: 'category_id' })
  category!: Category;

  @Property({ fieldName: 'image_url', type: 'text', nullable: true })
  imageUrl?: string;

  @OneToMany(() => MarketplaceProduct, mp => mp.product)
  marketplaceProducts = new Collection<MarketplaceProduct>(this);

  @ManyToMany(() => ProductAttribute, attr => attr.products, { owner: true, pivotTable: 'product_attributes' })
  attributes = new Collection<ProductAttribute>(this);
}
