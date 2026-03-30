/**
 * @jest-environment node
 */
import { Test, TestingModule } from '@nestjs/testing';
import { EntityManager } from '@mikro-orm/core';
import { PricesService } from './prices.service';
import { Price } from '@sales-monitor/database';

describe('PricesService', () => {
  let service: PricesService;
  let em: jest.Mocked<Pick<EntityManager, 'find'>>;

  const FIXED_NOW = new Date('2024-03-15T12:00:00.000Z');

  beforeEach(async () => {
    em = { find: jest.fn().mockResolvedValue([]) };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        PricesService,
        { provide: EntityManager, useValue: em },
      ],
    }).compile();

    service = module.get<PricesService>(PricesService);
    jest.useFakeTimers();
    jest.setSystemTime(FIXED_NOW);
  });

  afterEach(() => {
    jest.clearAllMocks();
    jest.useRealTimers();
  });

  describe('getLatestPrices', () => {
    it('uses default limit of 100', async () => {
      await service.getLatestPrices();

      expect(em.find).toHaveBeenCalledWith(
        Price,
        {},
        expect.objectContaining({ limit: 100 }),
      );
    });

    it('accepts a custom limit', async () => {
      await service.getLatestPrices(50);

      expect(em.find).toHaveBeenCalledWith(
        Price,
        {},
        expect.objectContaining({ limit: 50 }),
      );
    });

    it('orders by createdAt DESC', async () => {
      await service.getLatestPrices();

      expect(em.find).toHaveBeenCalledWith(
        Price,
        {},
        expect.objectContaining({ orderBy: { createdAt: 'DESC' } }),
      );
    });

    it('populates product brand and marketplace', async () => {
      await service.getLatestPrices();

      expect(em.find).toHaveBeenCalledWith(
        Price,
        {},
        expect.objectContaining({
          populate: ['marketplaceProduct.product.brand', 'marketplaceProduct.marketplace'],
        }),
      );
    });

    it('returns the array from em.find', async () => {
      const prices = [{ priceId: 1 }, { priceId: 2 }] as Price[];
      (em.find as jest.Mock).mockResolvedValue(prices);

      const result = await service.getLatestPrices();

      expect(result).toBe(prices);
    });
  });

  describe('getPriceTrends', () => {
    it('queries with correct dateFrom for default 30 days', async () => {
      await service.getPriceTrends();

      const expectedDateFrom = new Date(FIXED_NOW);
      expectedDateFrom.setDate(expectedDateFrom.getDate() - 30);

      expect(em.find).toHaveBeenCalledWith(
        Price,
        { createdAt: { $gte: expectedDateFrom } },
        expect.any(Object),
      );
    });

    it('queries with correct dateFrom for custom days=7', async () => {
      await service.getPriceTrends(undefined, 7);

      const expectedDateFrom = new Date(FIXED_NOW);
      expectedDateFrom.setDate(expectedDateFrom.getDate() - 7);

      expect(em.find).toHaveBeenCalledWith(
        Price,
        { createdAt: { $gte: expectedDateFrom } },
        expect.any(Object),
      );
    });

    it('does NOT include productId filter when not provided', async () => {
      await service.getPriceTrends();

      const [, where] = (em.find as jest.Mock).mock.calls[0];
      expect(where).not.toHaveProperty('marketplaceProduct.product.productId');
    });

    it('includes productId filter when provided', async () => {
      await service.getPriceTrends(5);

      const [, where] = (em.find as jest.Mock).mock.calls[0];
      expect(where['marketplaceProduct.product.productId']).toBe(5);
    });

    it('orders by createdAt ASC', async () => {
      await service.getPriceTrends();

      expect(em.find).toHaveBeenCalledWith(
        Price,
        expect.any(Object),
        expect.objectContaining({ orderBy: { createdAt: 'ASC' } }),
      );
    });

    it('populates product and marketplace', async () => {
      await service.getPriceTrends();

      expect(em.find).toHaveBeenCalledWith(
        Price,
        expect.any(Object),
        expect.objectContaining({
          populate: ['marketplaceProduct.product', 'marketplaceProduct.marketplace'],
        }),
      );
    });

    it('returns the array from em.find', async () => {
      const prices = [{ priceId: 1 }] as Price[];
      (em.find as jest.Mock).mockResolvedValue(prices);

      const result = await service.getPriceTrends();

      expect(result).toBe(prices);
    });
  });
});
