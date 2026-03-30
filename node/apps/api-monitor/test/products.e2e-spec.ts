import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, ValidationPipe } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import * as request from 'supertest';

import { ProductsController } from '../src/products/products.controller';
import { ProductsService } from '../src/products/products.service';
import { Product } from '@sales-monitor/database';

const makeCollection = (items: any[] = []) => ({
  getItems: jest.fn().mockReturnValue(items),
  contains: jest.fn().mockReturnValue(false),
  add: jest.fn(),
  remove: jest.fn(),
});

const mockEm = {
  findOne: jest.fn(),
  findAndCount: jest.fn(),
};

const mockProduct = {
  productId: 1,
  name: 'Test Shampoo',
  imageUrl: null,
  brand: { brandId: 1, name: 'Head & Shoulders' },
  category: { categoryId: 1, name: 'Shampoos' },
  attributes: makeCollection([]),
  marketplaceProducts: makeCollection([]),
} as unknown as Product;

describe('Products (e2e)', () => {
  let app: INestApplication;

  beforeAll(async () => {
    const moduleRef: TestingModule = await Test.createTestingModule({
      controllers: [ProductsController],
      providers: [
        ProductsService,
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

  // ── GET /api/products ─────────────────────────────────────

  describe('GET /api/products', () => {
    it('returns 200 with data and meta', async () => {
      (mockEm.findAndCount as jest.Mock).mockResolvedValue([[mockProduct], 1]);

      const res = await request(app.getHttpServer()).get('/api/products');

      expect(res.status).toBe(200);
      expect(res.body).toHaveProperty('data');
      expect(res.body).toHaveProperty('meta');
      expect(Array.isArray(res.body.data)).toBe(true);
    });

    it('passes categoryId filter to service', async () => {
      (mockEm.findAndCount as jest.Mock).mockResolvedValue([[], 0]);

      await request(app.getHttpServer()).get('/api/products?categoryId=3');

      const [, where] = (mockEm.findAndCount as jest.Mock).mock.calls[0];
      expect(where).toEqual({ category: { categoryId: 3 } });
    });

    it('passes brandId filter to service', async () => {
      (mockEm.findAndCount as jest.Mock).mockResolvedValue([[], 0]);

      await request(app.getHttpServer()).get('/api/products?brandId=7');

      const [, where] = (mockEm.findAndCount as jest.Mock).mock.calls[0];
      expect(where).toEqual({ brand: { brandId: 7 } });
    });

    it('passes search filter to service', async () => {
      (mockEm.findAndCount as jest.Mock).mockResolvedValue([[], 0]);

      await request(app.getHttpServer()).get('/api/products?search=shampoo');

      const [, where] = (mockEm.findAndCount as jest.Mock).mock.calls[0];
      expect(where).toEqual({ name: { $like: '%shampoo%' } });
    });

    it('applies page and limit to offset calculation', async () => {
      (mockEm.findAndCount as jest.Mock).mockResolvedValue([[], 50]);

      await request(app.getHttpServer()).get('/api/products?page=2&limit=10');

      const [, , opts] = (mockEm.findAndCount as jest.Mock).mock.calls[0];
      expect(opts.offset).toBe(10);
      expect(opts.limit).toBe(10);
    });
  });

  // ── GET /api/products/:id ─────────────────────────────────

  describe('GET /api/products/:id', () => {
    it('returns 200 with product when found', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockProduct);

      const res = await request(app.getHttpServer()).get('/api/products/1');

      expect(res.status).toBe(200);
      expect(res.body.productId).toBe(1);
    });

    it('returns 404 when product not found', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(null);

      const res = await request(app.getHttpServer()).get('/api/products/999');

      expect(res.status).toBe(404);
    });

    it('returns 400 for non-numeric id', async () => {
      const res = await request(app.getHttpServer()).get('/api/products/abc');

      expect(res.status).toBe(400);
    });
  });

  // ── GET /api/products/:id/prices ──────────────────────────

  describe('GET /api/products/:id/prices', () => {
    it('returns 200 with price history when product found', async () => {
      const productWithPrices = {
        ...mockProduct,
        marketplaceProducts: makeCollection([
          {
            marketplace: { marketplaceId: 1, name: 'Rozetka' },
            url: 'https://rozetka.com/p/1',
            prices: makeCollection([
              { priceId: 1, regularPrice: 100, createdAt: new Date('2024-01-01') },
            ]),
          },
        ]),
      } as unknown as Product;
      (mockEm.findOne as jest.Mock).mockResolvedValue(productWithPrices);

      const res = await request(app.getHttpServer()).get('/api/products/1/prices');

      expect(res.status).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
    });

    it('returns 404 when product not found', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(null);

      const res = await request(app.getHttpServer()).get('/api/products/999/prices');

      expect(res.status).toBe(404);
    });
  });
});
