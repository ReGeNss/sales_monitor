import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, ValidationPipe } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import * as request from 'supertest';

import { BrandsController } from '../src/brands/brands.controller';
import { BrandsService } from '../src/brands/brands.service';
import { Brand } from '@sales-monitor/database';

const makeCollection = (items: any[] = []) => ({
  getItems: jest.fn().mockReturnValue(items),
});

const mockEm = {
  find: jest.fn(),
  findOne: jest.fn(),
};

const mockBrand = {
  brandId: 1,
  name: 'Nike',
  bannerUrl: null,
  products: makeCollection([
    { productId: 1, name: 'Shoes', category: { categoryId: 1, name: 'Footwear' } },
    { productId: 2, name: 'T-Shirt', category: { categoryId: 2, name: 'Clothing' } },
    { productId: 3, name: 'Cap', category: { categoryId: 2, name: 'Clothing' } },
  ]),
} as unknown as Brand;

describe('Brands (e2e)', () => {
  let app: INestApplication;

  beforeAll(async () => {
    const moduleRef: TestingModule = await Test.createTestingModule({
      controllers: [BrandsController],
      providers: [
        BrandsService,
        { provide: EntityManager, useValue: mockEm },
        { provide: 'APP_GUARD', useValue: { canActivate: () => true } },
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
  beforeEach(() => {
    (mockEm.find as jest.Mock).mockResolvedValue([]);
    (mockEm.findOne as jest.Mock).mockResolvedValue(null);
  });
  afterEach(() => jest.clearAllMocks());

  // ── GET /api/brands ───────────────────────────────────────

  describe('GET /api/brands', () => {
    it('returns 200 with array of brands', async () => {
      (mockEm.find as jest.Mock).mockResolvedValue([mockBrand]);

      const res = await request(app.getHttpServer()).get('/api/brands');

      expect(res.status).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
      expect(res.body[0].brandId).toBe(1);
    });
  });

  // ── GET /api/brands/:id ───────────────────────────────────

  describe('GET /api/brands/:id', () => {
    it('returns 200 with brand when found', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockBrand);

      const res = await request(app.getHttpServer()).get('/api/brands/1');

      expect(res.status).toBe(200);
      expect(res.body.brandId).toBe(1);
    });

    it('returns 404 when brand not found', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(null);

      const res = await request(app.getHttpServer()).get('/api/brands/999');

      expect(res.status).toBe(404);
    });
  });

  // ── GET /api/brands/:id/products ──────────────────────────

  describe('GET /api/brands/:id/products', () => {
    it('returns 200 with paginated products', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockBrand);

      const res = await request(app.getHttpServer()).get('/api/brands/1/products?page=1&limit=20');

      expect(res.status).toBe(200);
      expect(res.body).toHaveProperty('data');
      expect(res.body).toHaveProperty('meta');
      expect(Array.isArray(res.body.data)).toBe(true);
    });

    it('returns 404 when brand not found', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(null);

      const res = await request(app.getHttpServer()).get('/api/brands/999/products?page=1&limit=20');

      expect(res.status).toBe(404);
    });

    it('returns correct pagination meta with custom page and limit', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockBrand);

      const res = await request(app.getHttpServer()).get('/api/brands/1/products?page=1&limit=2');

      expect(res.status).toBe(200);
      expect(res.body.meta.page).toBe(1);
      expect(res.body.meta.limit).toBe(2);
      expect(res.body.meta.total).toBe(3);
      expect(res.body.data).toHaveLength(2);
    });
  });
});
