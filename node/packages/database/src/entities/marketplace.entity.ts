import { Entity, PrimaryKey, Property, OneToMany, Collection } from '@mikro-orm/core';
import { MarketplaceProduct } from './marketplace-product.entity';

@Entity({ tableName: 'marketplaces' })
export class Marketplace {
  @PrimaryKey({ fieldName: 'marketplace_id' })
  marketplaceId!: number;

  @Property({ fieldName: 'name', length: 255, unique: true })
  name!: string;

  @Property({ fieldName: 'url', type: 'text' })
  url!: string;

  @OneToMany(() => MarketplaceProduct, mp => mp.marketplace)
  marketplaceProducts = new Collection<MarketplaceProduct>(this);
}
