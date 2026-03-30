import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, ValidationPipe } from '@nestjs/common';
import { APP_GUARD } from '@nestjs/core';
import { JwtModule, JwtService } from '@nestjs/jwt';
import { PassportModule } from '@nestjs/passport';
import { Reflector } from '@nestjs/core';
import { EntityManager } from '@mikro-orm/core';
import * as request from 'supertest';

import { UsersController } from '../src/users/users.controller';
import { UsersService } from '../src/users/users.service';
import { JwtAuthGuard } from '../src/auth/guards/jwt-auth.guard';
import { JwtStrategy } from '../src/auth/strategies/jwt.strategy';
import { RegisteredUserGuard } from '../src/auth/guards/registered-user.guard';
import { User } from '@sales-monitor/database';

const JWT_SECRET = 'your-secret-key-change-in-production';

const mockEm = {
  findOne: jest.fn(),
  flush: jest.fn().mockResolvedValue(undefined),
  removeAndFlush: jest.fn().mockResolvedValue(undefined),
};

const mockUser = {
  userId: 1,
  login: 'testuser',
  password: 'hashed',
  nfToken: undefined as string | undefined,
} as unknown as User;

describe('Users (e2e)', () => {
  let app: INestApplication;
  let registeredToken: string;
  let guestToken: string;

  beforeAll(async () => {
    const jwtService = new JwtService({ secret: JWT_SECRET });
    registeredToken = jwtService.sign({ sub: 1, login: 'testuser' });
    guestToken = jwtService.sign({
      isGuest: true,
      deviceModel: 'iPhone',
      appVersion: '1.0',
      platform: 'ios',
      locale: 'en-US',
    });

    const moduleRef: TestingModule = await Test.createTestingModule({
      imports: [
        PassportModule,
        JwtModule.register({ secret: JWT_SECRET, signOptions: { expiresIn: '24h' } }),
      ],
      controllers: [UsersController],
      providers: [
        UsersService,
        JwtStrategy,
        Reflector,
        RegisteredUserGuard,
        { provide: APP_GUARD, useClass: JwtAuthGuard },
        { provide: EntityManager, useValue: mockEm },
      ],
    }).compile();

    app = moduleRef.createNestApplication();
    app.setGlobalPrefix('api');
    app.useGlobalPipes(
      new ValidationPipe({
        whitelist: true,
        transform: true,
        transformOptions: { enableImplicitConversion: true },
      }),
    );
    await app.init();
  });

  afterAll(() => app.close());
  afterEach(() => jest.clearAllMocks());

  // ── PUT /api/users/notification-token ─────────────────────

  describe('PUT /api/users/notification-token', () => {
    it('returns 401 without token', async () => {
      const res = await request(app.getHttpServer())
        .put('/api/users/notification-token')
        .send({ nfToken: 'abc' });

      expect(res.status).toBe(401);
    });

    it('returns 403 with guest token', async () => {
      const res = await request(app.getHttpServer())
        .put('/api/users/notification-token')
        .set('Authorization', `Bearer ${guestToken}`)
        .send({ nfToken: 'abc' });

      expect(res.status).toBe(403);
    });

    it('returns 200 with registered token', async () => {
      const user = { ...mockUser } as unknown as User;
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(user)   // JwtStrategy.validate
        .mockResolvedValueOnce(user);  // UsersService

      const res = await request(app.getHttpServer())
        .put('/api/users/notification-token')
        .set('Authorization', `Bearer ${registeredToken}`)
        .send({ nfToken: 'firebase-token-abc' });

      expect(res.status).toBe(200);
      expect(res.body.message).toBe('Notification token updated successfully');
    });

    it('returns 200 with empty body (optional field)', async () => {
      const user = { ...mockUser } as unknown as User;
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(user)
        .mockResolvedValueOnce(user);

      const res = await request(app.getHttpServer())
        .put('/api/users/notification-token')
        .set('Authorization', `Bearer ${registeredToken}`)
        .send({});

      expect(res.status).toBe(200);
    });

    it('returns 404 when user not found in service', async () => {
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)  // JwtStrategy.validate
        .mockResolvedValueOnce(null);     // UsersService

      const res = await request(app.getHttpServer())
        .put('/api/users/notification-token')
        .set('Authorization', `Bearer ${registeredToken}`)
        .send({ nfToken: 'token' });

      expect(res.status).toBe(404);
    });
  });

  // ── DELETE /api/users/account ─────────────────────────────

  describe('DELETE /api/users/account', () => {
    it('returns 401 without token', async () => {
      const res = await request(app.getHttpServer()).delete('/api/users/account');

      expect(res.status).toBe(401);
    });

    it('returns 403 with guest token', async () => {
      const res = await request(app.getHttpServer())
        .delete('/api/users/account')
        .set('Authorization', `Bearer ${guestToken}`);

      expect(res.status).toBe(403);
    });

    it('returns 200 with registered token', async () => {
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)  // JwtStrategy.validate
        .mockResolvedValueOnce(mockUser); // UsersService

      const res = await request(app.getHttpServer())
        .delete('/api/users/account')
        .set('Authorization', `Bearer ${registeredToken}`);

      expect(res.status).toBe(200);
      expect(res.body.message).toBe('Account deleted successfully');
    });

    it('returns 404 when user not found in service', async () => {
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)  // JwtStrategy.validate
        .mockResolvedValueOnce(null);     // UsersService

      const res = await request(app.getHttpServer())
        .delete('/api/users/account')
        .set('Authorization', `Bearer ${registeredToken}`);

      expect(res.status).toBe(404);
    });
  });
});
