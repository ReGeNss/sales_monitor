export class MarketplaceDomain {
  constructor(
    public readonly marketplaceId: number,
    public readonly name: string,
    public readonly url: string,
  ) {}
}
