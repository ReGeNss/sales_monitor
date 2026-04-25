export class ProductAttributeDomain {
  constructor(
    public readonly attributeId: number,
    public readonly attributeType: string,
    public readonly value: string,
  ) {}
}
