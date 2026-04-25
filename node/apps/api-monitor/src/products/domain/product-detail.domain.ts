import { BrandDomain } from '../../brands/domain/brand.domain';
import { CategoryDomain } from '../../categories/domain/category.domain';
import { ProductAttributeDomain } from './product-attribute.domain';
import { MarketplaceProductDomain } from './marketplace-product.domain';

export class ProductDetailDomain {
  constructor(
    public readonly productId: number,
    public readonly name: string,
    public readonly imageUrl: string | undefined,
    public readonly brand: BrandDomain,
    public readonly category: CategoryDomain,
    public readonly attributes: ProductAttributeDomain[],
    public readonly marketplaceProducts: MarketplaceProductDomain[],
  ) {}
}
