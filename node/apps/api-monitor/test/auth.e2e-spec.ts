import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, ValidationPipe } from '@nestjs/common';
import { JwtModule } from '@nestjs/jwt';
import { PassportModule } from '@nestjs/passport';
import { EntityManager } from '@mikro-orm/core';
import * as request from 'supertest';
import * as bcrypt from 'bcrypt';

import { AuthController } from '../src/auth/auth.controller';
import { AuthService } from '../src/auth/auth.service';
import { LocalStrategy } from '../src/auth/strategies/local.strategy';
import { User } from '@sales-monitor/database';

jest.mock('bcrypt');

const JWT_SECRET = 'your-secret-key-change-in-production';

const mockEm = {
  findOne: jest.fn(),
  persistAndFlush: jest.fn().mockResolvedValue(undefined),
};

const mockUser = {
  userId: 1,
  login: 'testuser',
  password: '$2b$10$hashedPassword',
} as unknown as User;

describe('Auth (e2e)', () => {
  let app: INestApplication;

  beforeAll(async () => {
    const moduleRef: TestingModule = await Test.createTestingModule({
      imports: [
        PassportModule,
        JwtModule.register({ secret: JWT_SECRET, signOptions: { expiresIn: '24h' } }),
      ],
      controllers: [AuthController],
      providers: [
        AuthService,
        LocalStrategy,
        { provide: EntityManager, useValue: mockEm },
        { provide: 'APP_GUARD', useValue: { canActivate: () => true } },
      ],
    }).compile();

    app = moduleRef.createNestApplication();
    app.setGlobalPrefix('api');
    app.useGlobalPipes(
      new ValidationPipe({
        whitelist: true,
        forbidNonWhitelisted: true,
        transform: true,
        transformOptions: { enableImplicitConversion: true },
      }),
    );
    await app.init();
  });

  afterAll(() => app.close());
  afterEach(() => jest.clearAllMocks());

  // ── Guest Login ───────────────────────────────────────────

  describe('POST /api/auth/guest', () => {
    const validBody = {
      deviceModel: 'iPhone 15',
      appVersion: '1.0.0',
      platform: 'ios',
      locale: 'uk-UA',
    };

    it('returns 201 with access_token for valid body', async () => {
      const res = await request(app.getHttpServer())
        .post('/api/auth/guest')
        .send(validBody);

      expect(res.status).toBe(201);
      expect(res.body).toHaveProperty('access_token');
      expect(typeof res.body.access_token).toBe('string');
    });

    it('returns 400 when required fields are missing', async () => {
      const res = await request(app.getHttpServer())
        .post('/api/auth/guest')
        .send({ deviceModel: 'iPhone' });

      expect(res.status).toBe(400);
    });

    it('returns 400 for unknown extra fields', async () => {
      const res = await request(app.getHttpServer())
        .post('/api/auth/guest')
        .send({ ...validBody, unknownField: 'value' });

      expect(res.status).toBe(400);
    });
  });

  // ── Register ──────────────────────────────────────────────

  describe('POST /api/auth/register', () => {
    it('returns 201 with access_token and user for new login', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(null);

      const res = await request(app.getHttpServer())
        .post('/api/auth/register')
        .send({ login: 'newuser', password: 'password123' });

      expect(res.status).toBe(201);
      expect(res.body).toHaveProperty('access_token');
      expect(res.body.user.login).toBe('newuser');
    });

    it('returns 409 when login already exists', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockUser);

      const res = await request(app.getHttpServer())
        .post('/api/auth/register')
        .send({ login: 'testuser', password: 'password123' });

      expect(res.status).toBe(409);
    });

    it('returns 400 when password is too short', async () => {
      const res = await request(app.getHttpServer())
        .post('/api/auth/register')
        .send({ login: 'someuser', password: '123' });

      expect(res.status).toBe(400);
    });

    it('returns 400 when login is too short', async () => {
      const res = await request(app.getHttpServer())
        .post('/api/auth/register')
        .send({ login: 'ab', password: 'password123' });

      expect(res.status).toBe(400);
    });

    it('returns 400 when login is missing', async () => {
      const res = await request(app.getHttpServer())
        .post('/api/auth/register')
        .send({ password: 'password123' });

      expect(res.status).toBe(400);
    });
  });

  // ── Login ─────────────────────────────────────────────────

  describe('POST /api/auth/login', () => {
    it('returns 201 with access_token for valid credentials', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockUser);
      (bcrypt.compare as jest.Mock).mockResolvedValue(true);

      const res = await request(app.getHttpServer())
        .post('/api/auth/login')
        .send({ login: 'testuser', password: 'password123' });

      expect(res.status).toBe(201);
      expect(res.body).toHaveProperty('access_token');
      expect(res.body.user.login).toBe('testuser');
    });

    it('returns 401 when password is wrong', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockUser);
      (bcrypt.compare as jest.Mock).mockResolvedValue(false);

      const res = await request(app.getHttpServer())
        .post('/api/auth/login')
        .send({ login: 'testuser', password: 'wrongpassword' });

      expect(res.status).toBe(401);
    });

    it('returns 401 when user does not exist', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(null);

      const res = await request(app.getHttpServer())
        .post('/api/auth/login')
        .send({ login: 'ghost', password: 'password123' });

      expect(res.status).toBe(401);
    });
  });
});
