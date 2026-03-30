import { Test, TestingModule } from '@nestjs/testing';
import { NotFoundException } from '@nestjs/common';
import { EntityManager } from '@mikro-orm/core';
import { FavoritesService } from './favorites.service';
import { User, Product, Brand } from '@sales-monitor/database';

const makeCollection = (items: any[] = [], contains = false) => ({
  getItems: jest.fn().mockReturnValue(items),
  contains: jest.fn().mockReturnValue(contains),
  add: jest.fn(),
  remove: jest.fn(),
});

describe('FavoritesService', () => {
  let service: FavoritesService;
  let em: jest.Mocked<Pick<EntityManager, 'findOne' | 'flush'>>;

  const mockProduct = { productId: 1, name: 'Shampoo' } as unknown as Product;
  const mockBrand = { brandId: 1, name: 'Nike' } as unknown as Brand;

  const buildMockUser = (options: { containsProduct?: boolean; containsBrand?: boolean } = {}) =>
    ({
      userId: 1,
      login: 'testuser',
      favoriteProducts: makeCollection([mockProduct], options.containsProduct ?? false),
      favoriteBrands: makeCollection([mockBrand], options.containsBrand ?? false),
    }) as unknown as User;

  beforeEach(async () => {
    em = {
      findOne: jest.fn(),
      flush: jest.fn().mockResolvedValue(undefined),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        FavoritesService,
        { provide: EntityManager, useValue: em },
      ],
    }).compile();

    service = module.get<FavoritesService>(FavoritesService);
  });

  afterEach(() => jest.clearAllMocks());

  // ── Favorite Products ─────────────────────────────────────

  describe('getFavoriteProducts', () => {
    it('returns favorite products for the user', async () => {
      const mockUser = buildMockUser();
      (em.findOne as jest.Mock).mockResolvedValue(mockUser);

      const result = await service.getFavoriteProducts(1);

      expect(result).toEqual([mockProduct]);
    });

    it('throws NotFoundException when user not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(null);

      await expect(service.getFavoriteProducts(99)).rejects.toThrow(
        new NotFoundException('User not found'),
      );
    });
  });

  describe('addFavoriteProduct', () => {
    it('adds product to favorites when not already present', async () => {
      const mockUser = buildMockUser({ containsProduct: false });
      (em.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)
        .mockResolvedValueOnce(mockProduct);

      const result = await service.addFavoriteProduct(1, 1);

      expect(mockUser.favoriteProducts.add).toHaveBeenCalledWith(mockProduct);
      expect(em.flush).toHaveBeenCalledTimes(1);
      expect(result).toEqual({ message: 'Product added to favorites' });
    });

    it('does not add or flush when product already in favorites', async () => {
      const mockUser = buildMockUser({ containsProduct: true });
      (em.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)
        .mockResolvedValueOnce(mockProduct);

      await service.addFavoriteProduct(1, 1);

      expect(mockUser.favoriteProducts.add).not.toHaveBeenCalled();
      expect(em.flush).not.toHaveBeenCalled();
    });

    it('throws NotFoundException when user not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValueOnce(null);

      await expect(service.addFavoriteProduct(99, 1)).rejects.toThrow(
        new NotFoundException('User not found'),
      );
    });

    it('throws NotFoundException when product not found', async () => {
      const mockUser = buildMockUser();
      (em.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)
        .mockResolvedValueOnce(null);

      await expect(service.addFavoriteProduct(1, 999)).rejects.toThrow(
        new NotFoundException('Product not found'),
      );
    });
  });

  describe('removeFavoriteProduct', () => {
    it('removes product and flushes', async () => {
      const mockUser = buildMockUser();
      (em.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)
        .mockResolvedValueOnce(mockProduct);

      const result = await service.removeFavoriteProduct(1, 1);

      expect(mockUser.favoriteProducts.remove).toHaveBeenCalledWith(mockProduct);
      expect(em.flush).toHaveBeenCalledTimes(1);
      expect(result).toEqual({ message: 'Product removed from favorites' });
    });

    it('throws NotFoundException when user not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValueOnce(null);

      await expect(service.removeFavoriteProduct(99, 1)).rejects.toThrow(
        new NotFoundException('User not found'),
      );
    });

    it('throws NotFoundException when product not found', async () => {
      const mockUser = buildMockUser();
      (em.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)
        .mockResolvedValueOnce(null);

      await expect(service.removeFavoriteProduct(1, 999)).rejects.toThrow(
        new NotFoundException('Product not found'),
      );
    });
  });

  // ── Favorite Brands ───────────────────────────────────────

  describe('getFavoriteBrands', () => {
    it('returns favorite brands for the user', async () => {
      const mockUser = buildMockUser();
      (em.findOne as jest.Mock).mockResolvedValue(mockUser);

      const result = await service.getFavoriteBrands(1);

      expect(result).toEqual([mockBrand]);
    });

    it('throws NotFoundException when user not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValue(null);

      await expect(service.getFavoriteBrands(99)).rejects.toThrow(
        new NotFoundException('User not found'),
      );
    });
  });

  describe('addFavoriteBrand', () => {
    it('adds brand to favorites when not already present', async () => {
      const mockUser = buildMockUser({ containsBrand: false });
      (em.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)
        .mockResolvedValueOnce(mockBrand);

      const result = await service.addFavoriteBrand(1, 1);

      expect(mockUser.favoriteBrands.add).toHaveBeenCalledWith(mockBrand);
      expect(em.flush).toHaveBeenCalledTimes(1);
      expect(result).toEqual({ message: 'Brand added to favorites' });
    });

    it('does not add or flush when brand already in favorites', async () => {
      const mockUser = buildMockUser({ containsBrand: true });
      (em.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)
        .mockResolvedValueOnce(mockBrand);

      await service.addFavoriteBrand(1, 1);

      expect(mockUser.favoriteBrands.add).not.toHaveBeenCalled();
      expect(em.flush).not.toHaveBeenCalled();
    });

    it('throws NotFoundException when user not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValueOnce(null);

      await expect(service.addFavoriteBrand(99, 1)).rejects.toThrow(
        new NotFoundException('User not found'),
      );
    });

    it('throws NotFoundException when brand not found', async () => {
      const mockUser = buildMockUser();
      (em.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)
        .mockResolvedValueOnce(null);

      await expect(service.addFavoriteBrand(1, 999)).rejects.toThrow(
        new NotFoundException('Brand not found'),
      );
    });
  });

  describe('removeFavoriteBrand', () => {
    it('removes brand and flushes', async () => {
      const mockUser = buildMockUser();
      (em.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)
        .mockResolvedValueOnce(mockBrand);

      const result = await service.removeFavoriteBrand(1, 1);

      expect(mockUser.favoriteBrands.remove).toHaveBeenCalledWith(mockBrand);
      expect(em.flush).toHaveBeenCalledTimes(1);
      expect(result).toEqual({ message: 'Brand removed from favorites' });
    });

    it('throws NotFoundException when user not found', async () => {
      (em.findOne as jest.Mock).mockResolvedValueOnce(null);

      await expect(service.removeFavoriteBrand(99, 1)).rejects.toThrow(
        new NotFoundException('User not found'),
      );
    });

    it('throws NotFoundException when brand not found', async () => {
      const mockUser = buildMockUser();
      (em.findOne as jest.Mock)
        .mockResolvedValueOnce(mockUser)
        .mockResolvedValueOnce(null);

      await expect(service.removeFavoriteBrand(1, 999)).rejects.toThrow(
        new NotFoundException('Brand not found'),
      );
    });
  });
});
