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
exports.ProductAttribute = void 0;
const core_1 = require("@mikro-orm/core");
const product_entity_1 = require("./product.entity");
let ProductAttribute = class ProductAttribute {
    constructor() {
        this.products = new core_1.Collection(this);
    }
};
exports.ProductAttribute = ProductAttribute;
__decorate([
    (0, core_1.PrimaryKey)({ fieldName: 'attribute_id' }),
    __metadata("design:type", Number)
], ProductAttribute.prototype, "attributeId", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'name', length: 255 }),
    __metadata("design:type", String)
], ProductAttribute.prototype, "name", void 0);
__decorate([
    (0, core_1.Property)({ fieldName: 'value', length: 255 }),
    __metadata("design:type", String)
], ProductAttribute.prototype, "value", void 0);
__decorate([
    (0, core_1.ManyToMany)(() => product_entity_1.Product, product => product.attributes),
    __metadata("design:type", Object)
], ProductAttribute.prototype, "products", void 0);
exports.ProductAttribute = ProductAttribute = __decorate([
    (0, core_1.Entity)({ tableName: 'product_attributes' })
], ProductAttribute);
//# sourceMappingURL=product-attribute.entity.js.map