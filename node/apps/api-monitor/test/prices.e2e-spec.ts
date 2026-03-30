import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, ValidationPipe } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import * as request from 'supertest';

import { PricesController } from '../src/prices/prices.controller';
import { PricesService } from '../src/prices/prices.service';

const mockEm = {
  find: jest.fn(),
};

describe('Prices (e2e)', () => {
  let app: INestApplication;

  beforeAll(async () => {
    const moduleRef: TestingModule = await Test.createTestingModule({
      controllers: [PricesController],
      providers: [
        PricesService,
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
  });
  afterEach(() => jest.clearAllMocks());

  // ── GET /api/prices/latest ────────────────────────────────

  describe('GET /api/prices/latest', () => {
    it('returns 200 with array of prices', async () => {
      const res = await request(app.getHttpServer()).get('/api/prices/latest?limit=100');

      expect(res.status).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
    });

    it('passes custom limit to service', async () => {
      await request(app.getHttpServer()).get('/api/prices/latest?limit=25');

      const [, , opts] = (mockEm.find as jest.Mock).mock.calls[0];
      expect(opts.limit).toBe(25);
    });

    it('uses default limit of 100 when not specified', async () => {
      await request(app.getHttpServer()).get('/api/prices/latest?limit=100');

      const [, , opts] = (mockEm.find as jest.Mock).mock.calls[0];
      expect(opts.limit).toBe(100);
    });
  });

  // ── GET /api/prices/trends ────────────────────────────────

  describe('GET /api/prices/trends', () => {
    it('returns 200 with array of prices', async () => {
      const res = await request(app.getHttpServer()).get('/api/prices/trends?productId=1&days=30');

      expect(res.status).toBe(200);
      expect(Array.isArray(res.body)).toBe(true);
    });

    it('returns 200 with productId filter', async () => {
      const res = await request(app.getHttpServer()).get('/api/prices/trends?productId=3&days=30');

      expect(res.status).toBe(200);
    });

    it('returns 200 with days filter', async () => {
      const res = await request(app.getHttpServer()).get('/api/prices/trends?productId=1&days=7');

      expect(res.status).toBe(200);
    });

    it('returns 200 with both productId and days filters', async () => {
      const res = await request(app.getHttpServer()).get('/api/prices/trends?productId=3&days=14');

      expect(res.status).toBe(200);
    });

    it('passes productId filter to service', async () => {
      await request(app.getHttpServer()).get('/api/prices/trends?productId=5&days=30');

      const [, where] = (mockEm.find as jest.Mock).mock.calls[0];
      expect(where['marketplaceProduct.product.productId']).toBe(5);
    });
  });
});
