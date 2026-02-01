import { Injectable, NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { User } from '@sales-monitor/database';

@Injectable()
export class UsersService {
  constructor(private readonly em: EntityManager) {}

  async getProfile(userId: number) {
    const user = await this.em.findOne(User, { userId });
    if (!user) {
      throw new NotFoundException('User not found');
    }
    return {
      userId: user.userId,
      login: user.login,
      nfToken: user.nfToken,
    };
  }

  async updateNotificationToken(userId: number, nfToken?: string) {
    const user = await this.em.findOne(User, { userId });
    if (!user) {
      throw new NotFoundException('User not found');
    }

    user.nfToken = nfToken;
    await this.em.flush();

    return { message: 'Notification token updated successfully' };
  }

  async deleteAccount(userId: number) {
    const user = await this.em.findOne(User, { userId });
    if (!user) {
      throw new NotFoundException('User not found');
    }

    await this.em.removeAndFlush(user);
    return { message: 'Account deleted successfully' };
  }
}
