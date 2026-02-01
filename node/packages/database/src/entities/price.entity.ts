import { Entity, PrimaryKey, Property, ManyToOne } from '@mikro-orm/core';
import { MarketplaceProduct } from './marketplace-product.entity';

@Entity({ tableName: 'prices' })
export class Price {
  @PrimaryKey({ fieldName: 'price_id' })
  priceId!: number;

  @ManyToOne(() => MarketplaceProduct, { fieldName: 'marketplace_product_id' })
  marketplaceProduct!: MarketplaceProduct;

  @Property({ fieldName: 'regular_price', type: 'decimal', precision: 10, scale: 2 })
  regularPrice!: number;

  @Property({ fieldName: 'discount_price', type: 'decimal', precision: 10, scale: 2, nullable: true })
  discountPrice?: number;

  @Property({ fieldName: 'created_at', defaultRaw: 'CURRENT_TIMESTAMP' })
  createdAt!: Date;
}
