"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Product = void 0;
const core_1 = require("@mikro-orm/core");
const brand_entity_1 = require("./brand.entity");
const category_entity_1 = require("./category.entity");
const marketplace_product_entity_1 = require("./marketplace-product.entity");
const product_attribute_entity_1 = require("./product-attribute.entity");
let Product = class Product {
    constructor() {
        this.marketplaceProducts = new core_1.Collection(this);
        this.attributes = new core_1.Collection(this);
    }
};
exports.Product = Product;
__decorate([
    (0, core_1.PrimaryKey)({ fieldName: 'product_id' }),
    __metadata("design:type", Number)
], Product.prototype, "productId", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'name_fingerprint', length: 255, nullable: true, unique: true }),
    __metadata("design:type", String)
], Product.prototype, "nameFingerprint", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'brand_id' }),
    __metadata("design:type", Number)
], Product.prototype, "brandId", void 0);
__decorate([
    (0, core_1.ManyToOne)(() => brand_entity_1.Brand, { fieldName: 'brand_id' }),
    __metadata("design:type", brand_entity_1.Brand)
], Product.prototype, "brand", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'name', length: 255 }),
    __metadata("design:type", String)
], Product.prototype, "name", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'category_id' }),
    __metadata("design:type", Number)
], Product.prototype, "categoryId", void 0);
__decorate([
    (0, core_1.ManyToOne)(() => category_entity_1.Category, { fieldName: 'category_id' }),
    __metadata("design:type", category_entity_1.Category)
], Product.prototype, "category", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'image_url', type: 'text', nullable: true }),
    __metadata("design:type", String)
], Product.prototype, "imageUrl", void 0);
__decorate([
    (0, core_1.OneToMany)(() => marketplace_product_entity_1.MarketplaceProduct, mp => mp.product),
    __metadata("design:type", Object)
], Product.prototype, "marketplaceProducts", void 0);
__decorate([
    (0, core_1.ManyToMany)(() => product_attribute_entity_1.ProductAttribute, attr => attr.products, { pivotTable: 'product_attributes' }),
    __metadata("design:type", Object)
], Product.prototype, "attributes", void 0);
exports.Product = Product = __decorate([
    (0, core_1.Entity)({ tableName: 'products' })
], Product);
//# sourceMappingURL=product.entity.js.map