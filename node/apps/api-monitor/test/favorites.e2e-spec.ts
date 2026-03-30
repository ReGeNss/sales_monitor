import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, ValidationPipe } from '@nestjs/common';
import { APP_GUARD } from '@nestjs/core';
import { JwtModule, JwtService } from '@nestjs/jwt';
import { PassportModule } from '@nestjs/passport';
import { Reflector } from '@nestjs/core';
import { EntityManager } from '@mikro-orm/core';
import * as request from 'supertest';

import { FavoritesController } from '../src/favorites/favorites.controller';
import { FavoritesService } from '../src/favorites/favorites.service';
import { JwtAuthGuard } from '../src/auth/guards/jwt-auth.guard';
import { JwtStrategy } from '../src/auth/strategies/jwt.strategy';
import { RegisteredUserGuard } from '../src/auth/guards/registered-user.guard';
import { User, Product, Brand } from '@sales-monitor/database';

const JWT_SECRET = 'your-secret-key-change-in-production';

const makeCollection = (items: any[] = [], contains = false) => ({
  getItems: jest.fn().mockReturnValue(items),
  contains: jest.fn().mockReturnValue(contains),
  add: jest.fn(),
  remove: jest.fn(),
});

const mockEm = {
  findOne: jest.fn(),
  flush: jest.fn().mockResolvedValue(undefined),
};

const mockProduct = { productId: 1, name: 'Shampoo' } as unknown as Product;
const mockBrand = { brandId: 1, name: 'Nike' } as unknown as Brand;

const buildMockUser = () =>
  ({
    userId: 1,
    login: 'testuser',
    favoriteProducts: makeCollection([mockProduct]),
    favoriteBrands: makeCollection([mockBrand]),
  }) as unknown as User;

describe('Favorites (e2e)', () => {
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
      controllers: [FavoritesController],
      providers: [
        FavoritesService,
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

  // ── Guard enforcement ─────────────────────────────────────

  describe('Guard enforcement — all routes', () => {
    const routes = [
      { method: 'get', path: '/api/favorites/products' },
      { method: 'post', path: '/api/favorites/products/1' },
      { method: 'delete', path: '/api/favorites/products/1' },
      { method: 'get', path: '/api/favorites/brands' },
      { method: 'post', path: '/api/favorites/brands/1' },
      { method: 'delete', path: '/api/favorites/brands/1' },
    ] as const;

    routes.forEach(({ method, path }) => {
      it(`${method.toUpperCase()} ${path} returns 401 without token`, async () => {
        const res = await (request(app.getHttpServer()) as any)[method](path);
        expect(res.status).toBe(401);
      });

      it(`${method.toUpperCase()} ${path} returns 403 with guest token`, async () => {
        const res = await (request(app.getHttpServer()) as any)
          [method](path)
          .set('Authorization', `Bearer ${guestToken}`);
        expect(res.status).toBe(403);
      });
    });
  });

  // ── Favorite Products ─────────────────────────────────────

  describe('GET /api/favorites/products', () => {
    it('returns 200 with list of favorite products', async () => {
      const mockUser = buildMockUser();
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockUser);

      const res = await request(app.getHttpServer())
        .get('/api/favorites/products')
        .set('Authorization', `Bearer ${registeredToken}`);

      expect(res.status).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
    });
  });

  describe('POST /api/favorites/products/:productId', () => {
    it('returns 201 when product added successfully', async () => {
      const mockUser = buildMockUser();
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)    // JwtStrategy.validate
        .mockResolvedValueOnce(mockUser)    // FavoritesService: user lookup
        .mockResolvedValueOnce(mockProduct); // FavoritesService: product lookup

      const res = await request(app.getHttpServer())
        .post('/api/favorites/products/1')
        .set('Authorization', `Bearer ${registeredToken}`);

      expect(res.status).toBe(201);
      expect(res.body.message).toBe('Product added to favorites');
    });

    it('returns 404 when product not found', async () => {
      const mockUser = buildMockUser();
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)  // JwtStrategy.validate
        .mockResolvedValueOnce(mockUser)  // FavoritesService: user lookup
        .mockResolvedValueOnce(null);     // FavoritesService: product not found

      const res = await request(app.getHttpServer())
        .post('/api/favorites/products/999')
        .set('Authorization', `Bearer ${registeredToken}`);

      expect(res.status).toBe(404);
    });
  });

  describe('DELETE /api/favorites/products/:productId', () => {
    it('returns 200 when product removed successfully', async () => {
      const mockUser = buildMockUser();
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)    // JwtStrategy.validate
        .mockResolvedValueOnce(mockUser)    // FavoritesService: user lookup
        .mockResolvedValueOnce(mockProduct); // FavoritesService: product lookup

      const res = await request(app.getHttpServer())
        .delete('/api/favorites/products/1')
        .set('Authorization', `Bearer ${registeredToken}`);

      expect(res.status).toBe(200);
      expect(res.body.message).toBe('Product removed from favorites');
    });
  });

  // ── Favorite Brands ───────────────────────────────────────

  describe('GET /api/favorites/brands', () => {
    it('returns 200 with list of favorite brands', async () => {
      const mockUser = buildMockUser();
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockUser);

      const res = await request(app.getHttpServer())
        .get('/api/favorites/brands')
        .set('Authorization', `Bearer ${registeredToken}`);

      expect(res.status).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
    });
  });

  describe('POST /api/favorites/brands/:brandId', () => {
    it('returns 201 when brand added successfully', async () => {
      const mockUser = buildMockUser();
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)  // JwtStrategy.validate
        .mockResolvedValueOnce(mockUser)  // FavoritesService: user lookup
        .mockResolvedValueOnce(mockBrand); // FavoritesService: brand lookup

      const res = await request(app.getHttpServer())
        .post('/api/favorites/brands/1')
        .set('Authorization', `Bearer ${registeredToken}`);

      expect(res.status).toBe(201);
      expect(res.body.message).toBe('Brand added to favorites');
    });

    it('returns 404 when brand not found', async () => {
      const mockUser = buildMockUser();
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)  // JwtStrategy.validate
        .mockResolvedValueOnce(mockUser)  // FavoritesService: user lookup
        .mockResolvedValueOnce(null);     // FavoritesService: brand not found

      const res = await request(app.getHttpServer())
        .post('/api/favorites/brands/999')
        .set('Authorization', `Bearer ${registeredToken}`);

      expect(res.status).toBe(404);
    });
  });

  describe('DELETE /api/favorites/brands/:brandId', () => {
    it('returns 200 when brand removed successfully', async () => {
      const mockUser = buildMockUser();
      (mockEm.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)  // JwtStrategy.validate
        .mockResolvedValueOnce(mockUser)  // FavoritesService: user lookup
        .mockResolvedValueOnce(mockBrand); // FavoritesService: brand lookup

      const res = await request(app.getHttpServer())
        .delete('/api/favorites/brands/1')
        .set('Authorization', `Bearer ${registeredToken}`);

      expect(res.status).toBe(200);
      expect(res.body.message).toBe('Brand removed from favorites');
    });
  });
});
