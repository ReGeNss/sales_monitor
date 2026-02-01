import { Entity, PrimaryKey, Property, OneToMany, Collection } from '@mikro-orm/core';
import { Product } from './product.entity';

@Entity({ tableName: 'brands' })
export class Brand {
  @PrimaryKey({ fieldName: 'brand_id' })
  brandId!: number;

  @Property({ fieldName: 'name', length: 255, unique: true })
  name!: string;

  @Property({ fieldName: 'banner_url', type: 'text', nullable: true })
  bannerUrl?: string;

  @OneToMany(() => Product, product => product.brand)
  products = new Collection<Product>(this);
}
