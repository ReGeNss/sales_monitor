import { Injectable } from '@nestjs/common';
import { UsersRepository } from './users.repository';

@Injectable()
export class UsersService {
  constructor(private readonly usersRepository: UsersRepository) {}

  async updateNotificationToken(userId: number, nfToken?: string) {
    await this.usersRepository.updateNotificationToken(userId, nfToken);
    return { message: 'Notification token updated successfully' };
  }

  async deleteAccount(userId: number) {
    await this.usersRepository.delete(userId);
    return { message: 'Account deleted successfully' };
  }
}
