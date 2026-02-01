import { Test, TestingModule } from '@nestjs/testing';
import { JwtService } from '@nestjs/jwt';
import { EntityManager } from '@mikro-orm/core';
import { ConflictException, UnauthorizedException } from '@nestjs/common';
import { AuthService } from './auth.service';
import { User } from '@sales-monitor/database';
import * as bcrypt from 'bcrypt';

describe('AuthService', () => {
  let service: AuthService;
  let em: EntityManager;
  let jwtService: JwtService;

  const mockUser = {
    userId: 1,
    login: 'testuser',
    password: '$2b$10$hashedPassword',
  } as User;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        AuthService,
        {
          provide: EntityManager,
          useValue: {
            findOne: jest.fn(),
            create: jest.fn(),
            persistAndFlush: jest.fn(),
          },
        },
        {
          provide: JwtService,
          useValue: {
            sign: jest.fn().mockReturnValue('test-jwt-token'),
          },
        },
      ],
    }).compile();

    service = module.get<AuthService>(AuthService);
    em = module.get<EntityManager>(EntityManager);
    jwtService = module.get<JwtService>(JwtService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('register', () => {
    it('should register a new user successfully', async () => {
      jest.spyOn(em, 'findOne').mockResolvedValue(null);
      jest.spyOn(em, 'create').mockReturnValue(mockUser);
      jest.spyOn(em, 'persistAndFlush').mockResolvedValue(undefined);

      const result = await service.register({
        login: 'testuser',
        password: 'password123',
      });

      expect(result).toHaveProperty('access_token');
      expect(result).toHaveProperty('user');
      expect(result.user.login).toBe('testuser');
    });

    it('should throw ConflictException if user already exists', async () => {
      jest.spyOn(em, 'findOne').mockResolvedValue(mockUser);

      await expect(
        service.register({
          login: 'testuser',
          password: 'password123',
        }),
      ).rejects.toThrow(ConflictException);
    });
  });

  describe('validateUser', () => {
    it('should return user if credentials are valid', async () => {
      jest.spyOn(em, 'findOne').mockResolvedValue(mockUser);
      jest.spyOn(bcrypt, 'compare' as any).mockResolvedValue(true);

      const result = await service.validateUser('testuser', 'password123');

      expect(result).toEqual(mockUser);
    });

    it('should return null if user not found', async () => {
      jest.spyOn(em, 'findOne').mockResolvedValue(null);

      const result = await service.validateUser('testuser', 'password123');

      expect(result).toBeNull();
    });

    it('should return null if password is invalid', async () => {
      jest.spyOn(em, 'findOne').mockResolvedValue(mockUser);
      jest.spyOn(bcrypt, 'compare' as any).mockResolvedValue(false);

      const result = await service.validateUser('testuser', 'wrongpassword');

      expect(result).toBeNull();
    });
  });

  describe('login', () => {
    it('should return JWT token and user info', async () => {
      const result = await service.login(mockUser);

      expect(result).toHaveProperty('access_token');
      expect(result.access_token).toBe('test-jwt-token');
      expect(result.user.userId).toBe(1);
      expect(result.user.login).toBe('testuser');
    });
  });
});
