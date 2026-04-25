import { Injectable } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { User } from '@sales-monitor/database';
import { UserDomain } from '../common/domain/user.domain';
import { NotFoundError } from '../common/errors';

@Injectable()
export class UsersRepository {
  constructor(private readonly em: EntityManager) {}

  async findById(userId: number): Promise<UserDomain> {
    const user = await this.em.findOne(User, { userId });
    if (!user) {
      throw new NotFoundError('User not found');
    }
    return this.toDomain(user);
  }

  async updateNotificationToken(userId: number, nfToken?: string): Promise<void> {
    const user = await this.em.findOne(User, { userId });
    if (!user) {
      throw new NotFoundError('User not found');
    }
    user.nfToken = nfToken;
    await this.em.flush();
  }

  async delete(userId: number): Promise<void> {
    const user = await this.em.findOne(User, { userId });
    if (!user) {
      throw new NotFoundError('User not found');
    }
    await this.em.removeAndFlush(user);
  }

  private toDomain(orm: User): UserDomain {
    return new UserDomain(orm.userId, orm.login, orm.nfToken);
  }
}
