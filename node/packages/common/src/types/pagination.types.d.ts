export interface PaginationQuery {
    page?: number;
    limit?: number;
}
export interface PaginationMeta {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
    hasNext: boolean;
    hasPrev: boolean;
}
export interface PaginatedResponse<T> {
    data: T[];
    meta: PaginationMeta;
}
export declare function createPaginationMeta(page: number, limit: number, total: number): PaginationMeta;
