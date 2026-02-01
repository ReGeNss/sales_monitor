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
exports.Price = void 0;
const core_1 = require("@mikro-orm/core");
const marketplace_product_entity_1 = require("./marketplace-product.entity");
let Price = class Price {
};
exports.Price = Price;
__decorate([
    (0, core_1.PrimaryKey)({ fieldName: 'price_id' }),
    __metadata("design:type", Number)
], Price.prototype, "priceId", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'marketplace_product_id' }),
    __metadata("design:type", Number)
], Price.prototype, "marketplaceProductId", void 0);
__decorate([
    (0, core_1.ManyToOne)(() => marketplace_product_entity_1.MarketplaceProduct, { fieldName: 'marketplace_product_id' }),
    __metadata("design:type", marketplace_product_entity_1.MarketplaceProduct)
], Price.prototype, "marketplaceProduct", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'regular_price', type: 'decimal', precision: 10, scale: 2 }),
    __metadata("design:type", Number)
], Price.prototype, "regularPrice", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'discount_price', type: 'decimal', precision: 10, scale: 2, nullable: true }),
    __metadata("design:type", Number)
], Price.prototype, "discountPrice", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'created_at', defaultRaw: 'CURRENT_TIMESTAMP' }),
    __metadata("design:type", Date)
], Price.prototype, "createdAt", void 0);
exports.Price = Price = __decorate([
    (0, core_1.Entity)({ tableName: 'prices' })
], Price);
//# sourceMappingURL=price.entity.js.map