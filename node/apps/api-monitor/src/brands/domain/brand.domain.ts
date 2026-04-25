export class BrandDomain {
  constructor(
    public readonly brandId: number,
    public readonly name: string,
    public readonly bannerUrl?: string,
  ) {}
}
