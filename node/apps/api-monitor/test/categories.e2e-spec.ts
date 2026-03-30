import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, ValidationPipe } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import * as request from 'supertest';

import { CategoriesController } from '../src/categories/categories.controller';
import { CategoriesService } from '../src/categories/categories.service';
import { Category } from '@sales-monitor/database';

const makeCollection = (items: any[] = []) => ({
  getItems: jest.fn().mockReturnValue(items),
});

const mockEm = {
  find: jest.fn(),
  findOne: jest.fn(),
};

const mockCategory = {
  categoryId: 1,
  name: 'Electronics',
  products: makeCollection([
    { productId: 1, name: 'Phone', brand: { brandId: 1, name: 'Samsung' } },
    { productId: 2, name: 'Laptop', brand: { brandId: 1, name: 'Samsung' } },
    { productId: 3, name: 'Tablet', brand: { brandId: 2, name: 'Apple' } },
  ]),
} as unknown as Category;

describe('Categories (e2e)', () => {
  let app: INestApplication;

  beforeAll(async () => {
    const moduleRef: TestingModule = await Test.createTestingModule({
      controllers: [CategoriesController],
      providers: [
        CategoriesService,
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

  // ── GET /api/categories ───────────────────────────────────

  describe('GET /api/categories', () => {
    it('returns 200 with array of categories', async () => {
      (mockEm.find as jest.Mock).mockResolvedValue([mockCategory]);

      const res = await request(app.getHttpServer()).get('/api/categories');

      expect(res.status).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
      expect(res.body[0].categoryId).toBe(1);
    });
  });

  // ── GET /api/categories/:id ───────────────────────────────

  describe('GET /api/categories/:id', () => {
    it('returns 200 with category when found', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockCategory);

      const res = await request(app.getHttpServer()).get('/api/categories/1');

      expect(res.status).toBe(200);
      expect(res.body.categoryId).toBe(1);
    });

    it('returns 404 when category not found', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(null);

      const res = await request(app.getHttpServer()).get('/api/categories/999');

      expect(res.status).toBe(404);
    });
  });

  // ── GET /api/categories/:id/products ──────────────────────

  describe('GET /api/categories/:id/products', () => {
    it('returns 200 with paginated products', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockCategory);

      const res = await request(app.getHttpServer()).get('/api/categories/1/products?page=1&limit=20');

      expect(res.status).toBe(200);
      expect(res.body).toHaveProperty('data');
      expect(res.body).toHaveProperty('meta');
      expect(Array.isArray(res.body.data)).toBe(true);
    });

    it('returns 404 when category not found', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(null);

      const res = await request(app.getHttpServer()).get('/api/categories/999/products?page=1&limit=20');

      expect(res.status).toBe(404);
    });

    it('returns correct pagination meta with custom page and limit', async () => {
      (mockEm.findOne as jest.Mock).mockResolvedValue(mockCategory);

      const res = await request(app.getHttpServer()).get('/api/categories/1/products?page=1&limit=2');

      expect(res.status).toBe(200);
      expect(res.body.meta.page).toBe(1);
      expect(res.body.meta.limit).toBe(2);
      expect(res.body.meta.total).toBe(3);
      expect(res.body.data).toHaveLength(2);
    });
  });
});
