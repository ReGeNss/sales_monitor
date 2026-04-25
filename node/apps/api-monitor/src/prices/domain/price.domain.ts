export class PriceDomain {
  constructor(
    public readonly priceId: number,
    public readonly regularPrice: number,
    public readonly specialPrice: number | undefined,
    public readonly createdAt: Date,
  ) {}
}
