import { Test, TestingModule } from '@nestjs/testing';
import { NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { BrandsService } from './brands.service';
import { Brand } from '@sales-monitor/database';

const makeCollection = (items: any[] = []) => ({
  getItems: jest.fn().mockReturnValue(items),
});

describe('BrandsService', () => {
  let service: BrandsService;
  let em: jest.Mocked<Pick<EntityManager, 'find' | 'findOne'>>;

  beforeEach(async () => {
    em = {
      find: jest.fn(),
      findOne: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        BrandsService,
        { provide: EntityManager, useValue: em },
      ],
    }).compile();

    service = module.get<BrandsService>(BrandsService);
  });

  afterEach(() => jest.clearAllMocks());

  describe('findAll', () => {
    it('returns all brands ordered by name ASC', async () => {
      const brands = [
        { brandId: 2, name: 'Adidas' },
        { brandId: 1, name: 'Nike' },
      ] as Brand[];
      (em.find as jest.Mock).mockResolvedValue(brands);

      const result = await service.findAll();

      expect(em.find).toHaveBeenCalledWith(Brand, {}, { orderBy: { name: 'ASC' } });
      expect(result).toBe(brands);
    });
  });

  describe('findOne', () => {
    it('returns a brand when found', async () => {
      const brand = { brandId: 1, name: 'Nike' } as Brand;
      (em.findOne as jest.Mock).mockResolvedValue(brand);

      const result = await service.findOne(1);

      expect(result).toBe(brand);
    });

    it('throws NotFoundException when brand not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(null);

      await expect(service.findOne(5)).rejects.toThrow(
        new NotFoundException('Brand with ID 5 not found'),
      );
    });
  });

  describe('getBrandProducts', () => {
    it('throws NotFoundException when brand not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(null);

      await expect(service.getBrandProducts(99)).rejects.toThrow(
        new NotFoundException('Brand with ID 99 not found'),
      );
    });

    it('returns first page of products with correct meta', async () => {
      const products = [1, 2, 3, 4, 5].map((i) => ({ productId: i, name: `P${i}` }));
      const brand = { brandId: 1, name: 'Nike', products: makeCollection(products) } as unknown as Brand;
      (em.findOne as jest.Mock).mockResolvedValue(brand);

      const result = await service.getBrandProducts(1, 1, 3);

      expect(result.data).toHaveLength(3);
      expect(result.data).toEqual(products.slice(0, 3));
      expect(result.meta).toEqual({ page: 1, limit: 3, total: 5, totalPages: 2 });
    });

    it('returns second page of products', async () => {
      const products = [1, 2, 3, 4, 5].map((i) => ({ productId: i, name: `P${i}` }));
      const brand = { brandId: 1, name: 'Nike', products: makeCollection(products) } as unknown as Brand;
      (em.findOne as jest.Mock).mockResolvedValue(brand);

      const result = await service.getBrandProducts(1, 2, 3);

      expect(result.data).toHaveLength(2);
      expect(result.data).toEqual(products.slice(3));
      expect(result.meta.page).toBe(2);
      expect(result.meta.total).toBe(5);
      expect(result.meta.totalPages).toBe(2);
    });

    it('uses default page=1, limit=20', async () => {
      const products = Array.from({ length: 25 }, (_, i) => ({ productId: i + 1 }));
      const brand = { brandId: 1, name: 'Nike', products: makeCollection(products) } as unknown as Brand;
      (em.findOne as jest.Mock).mockResolvedValue(brand);

      const result = await service.getBrandProducts(1);

      expect(result.data).toHaveLength(20);
      expect(result.meta.limit).toBe(20);
      expect(result.meta.page).toBe(1);
    });

    it('calculates totalPages correctly', async () => {
      const products = Array.from({ length: 10 }, (_, i) => ({ productId: i + 1 }));
      const brand = { brandId: 1, name: 'Nike', products: makeCollection(products) } as unknown as Brand;
      (em.findOne as jest.Mock).mockResolvedValue(brand);

      const result = await service.getBrandProducts(1, 1, 3);

      expect(result.meta.totalPages).toBe(Math.ceil(10 / 3)); // 4
    });
  });
});
