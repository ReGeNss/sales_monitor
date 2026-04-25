import { Injectable } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { User } from '@sales-monitor/database';
import { UserDomain } from '../common/domain/user.domain';
import { ConflictError } from '../common/errors';

@Injectable()
export class AuthRepository {
  constructor(private readonly em: EntityManager) {}

  async findByLogin(login: string): Promise<{ user: UserDomain; hashedPassword: string } | null> {
    const user = await this.em.findOne(User, { login });
    if (!user) return null;
    return { user: this.toDomain(user), hashedPassword: user.password };
  }

  async createUser(login: string, hashedPassword: string): Promise<UserDomain> {
    const existing = await this.em.findOne(User, { login });
    if (existing) {
      throw new ConflictError('User with this login already exists');
    }
    const user = new User();
    user.login = login;
    user.password = hashedPassword;
    await this.em.persistAndFlush(user);
    return this.toDomain(user);
  }

  private toDomain(orm: User): UserDomain {
    return new UserDomain(orm.userId, orm.login, orm.nfToken);
  }
}
