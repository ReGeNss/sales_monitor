import { BrandDomain } from '../../brands/domain/brand.domain';
import { CategoryDomain } from '../../categories/domain/category.domain';

export class ProductDomain {
  constructor(
    public readonly productId: number,
    public readonly name: string,
    public readonly imageUrl: string | undefined,
    public readonly brand: BrandDomain,
    public readonly category: CategoryDomain,
  ) {}
}
