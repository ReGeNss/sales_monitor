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
exports.Marketplace = void 0;
const core_1 = require("@mikro-orm/core");
const marketplace_product_entity_1 = require("./marketplace-product.entity");
let Marketplace = class Marketplace {
    constructor() {
        this.marketplaceProducts = new core_1.Collection(this);
    }
};
exports.Marketplace = Marketplace;
__decorate([
    (0, core_1.PrimaryKey)({ fieldName: 'marketplace_id' }),
    __metadata("design:type", Number)
], Marketplace.prototype, "marketplaceId", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'name', length: 255, unique: true }),
    __metadata("design:type", String)
], Marketplace.prototype, "name", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'url', type: 'text' }),
    __metadata("design:type", String)
], Marketplace.prototype, "url", void 0);
__decorate([
    (0, core_1.OneToMany)(() => marketplace_product_entity_1.MarketplaceProduct, mp => mp.marketplace),
    __metadata("design:type", Object)
], Marketplace.prototype, "marketplaceProducts", void 0);
exports.Marketplace = Marketplace = __decorate([
    (0, core_1.Entity)({ tableName: 'marketplaces' })
], Marketplace);
//# sourceMappingURL=marketplace.entity.js.map