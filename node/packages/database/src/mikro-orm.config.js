"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const mysql_1 = require("@mikro-orm/mysql");
const entities = require("./entities");
exports.default = (0, mysql_1.defineConfig)({
    entities: Object.values(entities),
    host: process.env.DB_HOST || 'localhost',
    port: parseInt(process.env.DB_PORT || '3306'),
    dbName: process.env.DATABASE_NAME || 'sales_monitor',
    user: process.env.DB_USER_NAME || 'root',
    password: process.env.DB_USER_PASSWORD || '',
    debug: process.env.NODE_ENV === 'development',
    allowGlobalContext: true,
});
//# sourceMappingURL=mikro-orm.config.js.map