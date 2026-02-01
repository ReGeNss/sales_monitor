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
exports.MarketplaceProduct = void 0;
const core_1 = require("@mikro-orm/core");
const marketplace_entity_1 = require("./marketplace.entity");
const product_entity_1 = require("./product.entity");
const price_entity_1 = require("./price.entity");
let MarketplaceProduct = class MarketplaceProduct {
    constructor() {
        this.prices = new core_1.Collection(this);
    }
};
exports.MarketplaceProduct = MarketplaceProduct;
__decorate([
    (0, core_1.PrimaryKey)({ fieldName: 'marketplace_product_id' }),
    __metadata("design:type", Number)
], MarketplaceProduct.prototype, "marketplaceProductId", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'marketplace_id' }),
    __metadata("design:type", Number)
], MarketplaceProduct.prototype, "marketplaceId", void 0);
__decorate([
    (0, core_1.ManyToOne)(() => marketplace_entity_1.Marketplace, { fieldName: 'marketplace_id' }),
    __metadata("design:type", marketplace_entity_1.Marketplace)
], MarketplaceProduct.prototype, "marketplace", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'product_id' }),
    __metadata("design:type", Number)
], MarketplaceProduct.prototype, "productId", void 0);
__decorate([
    (0, core_1.ManyToOne)(() => product_entity_1.Product, { fieldName: 'product_id' }),
    __metadata("design:type", product_entity_1.Product)
], MarketplaceProduct.prototype, "product", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'url', type: 'text' }),
    __metadata("design:type", String)
], MarketplaceProduct.prototype, "url", void 0);
__decorate([
    (0, core_1.OneToMany)(() => price_entity_1.Price, price => price.marketplaceProduct),
    __metadata("design:type", Object)
], MarketplaceProduct.prototype, "prices", void 0);
exports.MarketplaceProduct = MarketplaceProduct = __decorate([
    (0, core_1.Entity)({ tableName: 'marketplace_products' })
], MarketplaceProduct);
//# sourceMappingURL=marketplace-product.entity.js.map