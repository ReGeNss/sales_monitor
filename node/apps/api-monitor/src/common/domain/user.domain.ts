export class UserDomain {
  constructor(
    public readonly userId: number,
    public readonly login: string,
    public readonly nfToken?: string,
  ) {}
}
