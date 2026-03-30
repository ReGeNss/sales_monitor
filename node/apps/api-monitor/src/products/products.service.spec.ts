import { Test, TestingModule } from '@nestjs/testing';
import { NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { ProductsService } from './products.service';
import { Product } from '@sales-monitor/database';

const makeCollection = (items: any[] = []) => ({
  getItems: jest.fn().mockReturnValue(items),
  contains: jest.fn().mockReturnValue(false),
  add: jest.fn(),
  remove: jest.fn(),
});

describe('ProductsService', () => {
  let service: ProductsService;
  let em: jest.Mocked<Pick<EntityManager, 'findOne' | 'findAndCount'>>;

  beforeEach(async () => {
    em = {
      findOne: jest.fn(),
      findAndCount: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ProductsService,
        { provide: EntityManager, useValue: em },
      ],
    }).compile();

    service = module.get<ProductsService>(ProductsService);
  });

  afterEach(() => jest.clearAllMocks());

  describe('findAll', () => {
    const mockProducts = [
      { productId: 2, name: 'Product B' },
      { productId: 1, name: 'Product A' },
    ] as Product[];

    it('returns paginated response with data and meta', async () => {
      (em.findAndCount as jest.Mock).mockResolvedValue([mockProducts, 2]);

      const result = await service.findAll({});

      expect(result.data).toBe(mockProducts);
      expect(result.meta).toEqual({
        page: 1,
        limit: 20,
        total: 2,
        totalPages: 1,
        hasNext: false,
        hasPrev: false,
      });
    });

    it('uses default page=1 and limit=20', async () => {
      (em.findAndCount as jest.Mock).mockResolvedValue([[], 0]);

      await service.findAll({});

      expect(em.findAndCount).toHaveBeenCalledWith(
        Product,
        {},
        expect.objectContaining({ limit: 20, offset: 0 }),
      );
    });

    it('applies categoryId filter', async () => {
      (em.findAndCount as jest.Mock).mockResolvedValue([[], 0]);

      await service.findAll({ categoryId: 3 });

      expect(em.findAndCount).toHaveBeenCalledWith(
        Product,
        { category: { categoryId: 3 } },
        expect.any(Object),
      );
    });

    it('applies brandId filter', async () => {
      (em.findAndCount as jest.Mock).mockResolvedValue([[], 0]);

      await service.findAll({ brandId: 7 });

      expect(em.findAndCount).toHaveBeenCalledWith(
        Product,
        { brand: { brandId: 7 } },
        expect.any(Object),
      );
    });

    it('applies search filter with $like', async () => {
      (em.findAndCount as jest.Mock).mockResolvedValue([[], 0]);

      await service.findAll({ search: 'shampoo' });

      expect(em.findAndCount).toHaveBeenCalledWith(
        Product,
        { name: { $like: '%shampoo%' } },
        expect.any(Object),
      );
    });

    it('calculates offset correctly for page 2', async () => {
      (em.findAndCount as jest.Mock).mockResolvedValue([[], 50]);

      await service.findAll({ page: 2, limit: 10 });

      expect(em.findAndCount).toHaveBeenCalledWith(
        Product,
        {},
        expect.objectContaining({ limit: 10, offset: 10 }),
      );
    });

    it('returns correct meta for multiple pages', async () => {
      (em.findAndCount as jest.Mock).mockResolvedValue([mockProducts, 50]);

      const result = await service.findAll({ page: 2, limit: 20 });

      expect(result.meta).toEqual({
        page: 2,
        limit: 20,
        total: 50,
        totalPages: 3,
        hasNext: true,
        hasPrev: true,
      });
    });

    it('orders by productId DESC', async () => {
      (em.findAndCount as jest.Mock).mockResolvedValue([[], 0]);

      await service.findAll({});

      expect(em.findAndCount).toHaveBeenCalledWith(
        Product,
        expect.any(Object),
        expect.objectContaining({ orderBy: { productId: 'DESC' } }),
      );
    });

    it('populates brand and category', async () => {
      (em.findAndCount as jest.Mock).mockResolvedValue([[], 0]);

      await service.findAll({});

      expect(em.findAndCount).toHaveBeenCalledWith(
        Product,
        expect.any(Object),
        expect.objectContaining({ populate: ['brand', 'category'] }),
      );
    });
  });

  describe('findOne', () => {
    it('returns product when found', async () => {
      const product = { productId: 1, name: 'Test' } as Product;
      (em.findOne as jest.Mock).mockResolvedValue(product);

      const result = await service.findOne(1);

      expect(result).toBe(product);
    });

    it('populates brand, category, attributes and marketplace', async () => {
      const product = { productId: 1 } as Product;
      (em.findOne as jest.Mock).mockResolvedValue(product);

      await service.findOne(1);

      expect(em.findOne).toHaveBeenCalledWith(
        Product,
        { productId: 1 },
        expect.objectContaining({
          populate: ['brand', 'category', 'attributes', 'marketplaceProducts.marketplace'],
        }),
      );
    });

    it('throws NotFoundException when product not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(null);

      await expect(service.findOne(42)).rejects.toThrow(
        new NotFoundException('Product with ID 42 not found'),
      );
    });
  });

  describe('getProductPrices', () => {
    it('throws NotFoundException when product not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(null);

      await expect(service.getProductPrices(99)).rejects.toThrow(
        new NotFoundException('Product with ID 99 not found'),
      );
    });

    it('returns mapped array with marketplace, url and prices', async () => {
      const marketplace = { marketplaceId: 1, name: 'Rozetka' };
      const prices = [
        { priceId: 1, regularPrice: 100, createdAt: new Date('2024-01-02') },
        { priceId: 2, regularPrice: 90, createdAt: new Date('2024-01-01') },
      ];
      const mp = {
        marketplace,
        url: 'https://rozetka.com/p/1',
        prices: makeCollection(prices),
      };
      const product = {
        productId: 1,
        marketplaceProducts: makeCollection([mp]),
      } as unknown as Product;
      (em.findOne as jest.Mock).mockResolvedValue(product);

      const result = await service.getProductPrices(1);

      expect(result).toHaveLength(1);
      expect(result[0].marketplace).toBe(marketplace);
      expect(result[0].url).toBe('https://rozetka.com/p/1');
    });

    it('sorts prices descending by createdAt', async () => {
      const prices = [
        { priceId: 1, createdAt: new Date('2024-01-01') },
        { priceId: 3, createdAt: new Date('2024-01-03') },
        { priceId: 2, createdAt: new Date('2024-01-02') },
      ];
      const mp = {
        marketplace: { marketplaceId: 1 },
        url: 'https://rozetka.com/p/1',
        prices: makeCollection(prices),
      };
      const product = {
        productId: 1,
        marketplaceProducts: makeCollection([mp]),
      } as unknown as Product;
      (em.findOne as jest.Mock).mockResolvedValue(product);

      const result = await service.getProductPrices(1);
      const returnedPrices = result[0].prices;

      expect(returnedPrices[0].priceId).toBe(3);
      expect(returnedPrices[1].priceId).toBe(2);
      expect(returnedPrices[2].priceId).toBe(1);
    });

    it('slices prices to a maximum of 30 records', async () => {
      const prices = Array.from({ length: 35 }, (_, i) => ({
        priceId: i + 1,
        createdAt: new Date(Date.now() - i * 1000),
      }));
      const mp = {
        marketplace: { marketplaceId: 1 },
        url: 'https://rozetka.com/p/1',
        prices: makeCollection(prices),
      };
      const product = {
        productId: 1,
        marketplaceProducts: makeCollection([mp]),
      } as unknown as Product;
      (em.findOne as jest.Mock).mockResolvedValue(product);

      const result = await service.getProductPrices(1);

      expect(result[0].prices).toHaveLength(30);
    });

    it('populates marketplaceProducts prices and marketplace', async () => {
      const mp = {
        marketplace: {},
        url: '',
        prices: makeCollection([]),
      };
      const product = {
        productId: 1,
        marketplaceProducts: makeCollection([mp]),
      } as unknown as Product;
      (em.findOne as jest.Mock).mockResolvedValue(product);

      await service.getProductPrices(1);

      expect(em.findOne).toHaveBeenCalledWith(
        Product,
        { productId: 1 },
        expect.objectContaining({
          populate: ['marketplaceProducts.prices', 'marketplaceProducts.marketplace'],
        }),
      );
    });
  });
});
