import { Entity, PrimaryKey, Property, ManyToMany, Collection } from '@mikro-orm/core';
import { Product } from './product.entity';

@Entity({ tableName: 'attributes' })
export class ProductAttribute {
  @PrimaryKey({ fieldName: 'attribute_id' })
  attributeId!: number;

  @Property({ fieldName: 'attribute_type' })
  attributeType!: string; // 'volume' or 'weight'

  @Property({ fieldName: 'value', type: 'text' })
  value!: string;

  @ManyToMany(() => Product, product => product.attributes, { mappedBy: 'attributes' })
  products = new Collection<Product>(this);
}
