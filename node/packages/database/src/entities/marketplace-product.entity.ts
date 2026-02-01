import { Entity, PrimaryKey, Property, ManyToOne, OneToMany, Collection } from '@mikro-orm/core';
import { Marketplace } from './marketplace.entity';
import { Product } from './product.entity';
import { Price } from './price.entity';

@Entity({ tableName: 'marketplace_products' })
export class MarketplaceProduct {
  @PrimaryKey({ fieldName: 'marketplace_product_id' })
  marketplaceProductId!: number;

  @ManyToOne(() => Marketplace, { fieldName: 'marketplace_id' })
  marketplace!: Marketplace;

  @ManyToOne(() => Product, { fieldName: 'product_id' })
  product!: Product;

  @Property({ fieldName: 'url', type: 'text' })
  url!: string;

  @OneToMany(() => Price, price => price.marketplaceProduct)
  prices = new Collection<Price>(this);
}
