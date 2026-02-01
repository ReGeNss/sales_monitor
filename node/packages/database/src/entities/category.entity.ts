import { Entity, PrimaryKey, Property, OneToMany, Collection } from '@mikro-orm/core';
import { Product } from './product.entity';

@Entity({ tableName: 'categories' })
export class Category {
  @PrimaryKey({ fieldName: 'category_id' })
  categoryId!: number;

  @Property({ fieldName: 'name', length: 255, unique: true })
  name!: string;

  @OneToMany(() => Product, product => product.category)
  products = new Collection<Product>(this);
}
