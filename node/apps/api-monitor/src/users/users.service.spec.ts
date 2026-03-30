import { Test, TestingModule } from '@nestjs/testing';
import { NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { UsersService } from './users.service';
import { User } from '@sales-monitor/database';

describe('UsersService', () => {
  let service: UsersService;
  let em: jest.Mocked<Pick<EntityManager, 'findOne' | 'flush' | 'removeAndFlush'>>;

  const mockUser = {
    userId: 1,
    login: 'testuser',
    nfToken: undefined as string | undefined,
  } as unknown as User;

  beforeEach(async () => {
    em = {
      findOne: jest.fn(),
      flush: jest.fn().mockResolvedValue(undefined),
      removeAndFlush: jest.fn().mockResolvedValue(undefined),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        UsersService,
        { provide: EntityManager, useValue: em },
      ],
    }).compile();

    service = module.get<UsersService>(UsersService);
  });

  afterEach(() => jest.clearAllMocks());

  describe('updateNotificationToken', () => {
    it('updates nfToken and returns success message', async () => {
      const user = { ...mockUser, nfToken: undefined as string | undefined } as unknown as User;
      (em.findOne as jest.Mock).mockResolvedValue(user);

      const result = await service.updateNotificationToken(1, 'firebase-token-xyz');

      expect(user.nfToken).toBe('firebase-token-xyz');
      expect(em.flush).toHaveBeenCalledTimes(1);
      expect(result).toEqual({ message: 'Notification token updated successfully' });
    });

    it('clears nfToken when called with undefined', async () => {
      const user = { ...mockUser, nfToken: 'old-token' } as unknown as User;
      (em.findOne as jest.Mock).mockResolvedValue(user);

      await service.updateNotificationToken(1, undefined);

      expect(user.nfToken).toBeUndefined();
      expect(em.flush).toHaveBeenCalledTimes(1);
    });

    it('throws NotFoundException when user not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(null);

      await expect(service.updateNotificationToken(99, 'token')).rejects.toThrow(
        new NotFoundException('User not found'),
      );
      expect(em.flush).not.toHaveBeenCalled();
    });
  });

  describe('deleteAccount', () => {
    it('removes user and returns success message', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(mockUser);

      const result = await service.deleteAccount(1);

      expect(em.removeAndFlush).toHaveBeenCalledWith(mockUser);
      expect(result).toEqual({ message: 'Account deleted successfully' });
    });

    it('throws NotFoundException when user not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(null);

      await expect(service.deleteAccount(99)).rejects.toThrow(
        new NotFoundException('User not found'),
      );
      expect(em.removeAndFlush).not.toHaveBeenCalled();
    });
  });
});
