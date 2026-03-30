import { Test, TestingModule } from '@nestjs/testing';
import { NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { CategoriesService } from './categories.service';
import { Category } from '@sales-monitor/database';

const makeCollection = (items: any[] = []) => ({
  getItems: jest.fn().mockReturnValue(items),
});

describe('CategoriesService', () => {
  let service: CategoriesService;
  let em: jest.Mocked<Pick<EntityManager, 'find' | 'findOne'>>;

  beforeEach(async () => {
    em = {
      find: jest.fn(),
      findOne: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        CategoriesService,
        { provide: EntityManager, useValue: em },
      ],
    }).compile();

    service = module.get<CategoriesService>(CategoriesService);
  });

  afterEach(() => jest.clearAllMocks());

  describe('findAll', () => {
    it('returns all categories ordered by name ASC', async () => {
      const categories = [
        { categoryId: 2, name: 'Electronics' },
        { categoryId: 1, name: 'Shampoos' },
      ] as Category[];
      (em.find as jest.Mock).mockResolvedValue(categories);

      const result = await service.findAll();

      expect(em.find).toHaveBeenCalledWith(Category, {}, { orderBy: { name: 'ASC' } });
      expect(result).toBe(categories);
    });
  });

  describe('findOne', () => {
    it('returns a category when found', async () => {
      const category = { categoryId: 1, name: 'Electronics' } as Category;
      (em.findOne as jest.Mock).mockResolvedValue(category);

      const result = await service.findOne(1);

      expect(result).toBe(category);
    });

    it('throws NotFoundException when category not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(null);

      await expect(service.findOne(7)).rejects.toThrow(
        new NotFoundException('Category with ID 7 not found'),
      );
    });
  });

  describe('getCategoryProducts', () => {
    it('throws NotFoundException when category not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(null);

      await expect(service.getCategoryProducts(99)).rejects.toThrow(
        new NotFoundException('Category with ID 99 not found'),
      );
    });

    it('returns first page of products with correct meta', async () => {
      const products = [1, 2, 3, 4, 5].map((i) => ({ productId: i, name: `P${i}` }));
      const category = {
        categoryId: 1,
        name: 'Electronics',
        products: makeCollection(products),
      } as unknown as Category;
      (em.findOne as jest.Mock).mockResolvedValue(category);

      const result = await service.getCategoryProducts(1, 1, 3);

      expect(result.data).toHaveLength(3);
      expect(result.data).toEqual(products.slice(0, 3));
      expect(result.meta).toEqual({ page: 1, limit: 3, total: 5, totalPages: 2 });
    });

    it('returns second page of products', async () => {
      const products = [1, 2, 3, 4, 5].map((i) => ({ productId: i, name: `P${i}` }));
      const category = {
        categoryId: 1,
        name: 'Electronics',
        products: makeCollection(products),
      } as unknown as Category;
      (em.findOne as jest.Mock).mockResolvedValue(category);

      const result = await service.getCategoryProducts(1, 2, 3);

      expect(result.data).toHaveLength(2);
      expect(result.data).toEqual(products.slice(3));
      expect(result.meta.page).toBe(2);
      expect(result.meta.totalPages).toBe(2);
    });

    it('uses default page=1, limit=20', async () => {
      const products = Array.from({ length: 25 }, (_, i) => ({ productId: i + 1 }));
      const category = {
        categoryId: 1,
        name: 'Electronics',
        products: makeCollection(products),
      } as unknown as Category;
      (em.findOne as jest.Mock).mockResolvedValue(category);

      const result = await service.getCategoryProducts(1);

      expect(result.data).toHaveLength(20);
      expect(result.meta.limit).toBe(20);
    });

    it('calculates totalPages correctly', async () => {
      const products = Array.from({ length: 10 }, (_, i) => ({ productId: i + 1 }));
      const category = {
        categoryId: 1,
        name: 'Electronics',
        products: makeCollection(products),
      } as unknown as Category;
      (em.findOne as jest.Mock).mockResolvedValue(category);

      const result = await service.getCategoryProducts(1, 1, 3);

      expect(result.meta.totalPages).toBe(Math.ceil(10 / 3)); // 4
    });
  });
});
