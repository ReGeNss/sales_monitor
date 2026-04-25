import { IDomainError } from './domain-error.interface';

export abstract class DomainError extends Error implements IDomainError {
  constructor(message: string) {
    super(message);
    this.name = this.constructor.name;
    Object.setPrototypeOf(this, new.target.prototype);
  }
}
