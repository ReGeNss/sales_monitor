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
exports.User = void 0;
const core_1 = require("@mikro-orm/core");
const product_entity_1 = require("./product.entity");
const brand_entity_1 = require("./brand.entity");
let User = class User {
    constructor() {
        this.favoriteProducts = new core_1.Collection(this);
        this.favoriteBrands = new core_1.Collection(this);
    }
};
exports.User = User;
__decorate([
    (0, core_1.PrimaryKey)({ fieldName: 'user_id', autoincrement: true }),
    __metadata("design:type", Number)
], User.prototype, "userId", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'login', length: 255, unique: true }),
    __metadata("design:type", String)
], User.prototype, "login", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'password', length: 255, hidden: true }),
    __metadata("design:type", String)
], User.prototype, "password", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'nf_token', type: 'text', nullable: true }),
    __metadata("design:type", String)
], User.prototype, "nfToken", void 0);
__decorate([
    (0, core_1.ManyToMany)(() => product_entity_1.Product, undefined, { pivotTable: 'favorite_products' }),
    __metadata("design:type", Object)
], User.prototype, "favoriteProducts", void 0);
__decorate([
    (0, core_1.ManyToMany)(() => brand_entity_1.Brand, undefined, { pivotTable: 'favorite_brands' }),
    __metadata("design:type", Object)
], User.prototype, "favoriteBrands", void 0);
exports.User = User = __decorate([
    (0, core_1.Entity)({ tableName: 'users' })
], User);
//# sourceMappingURL=user.entity.js.map