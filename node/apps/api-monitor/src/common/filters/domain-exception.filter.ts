import { ExceptionFilter, Catch, ArgumentsHost } from '@nestjs/common';
import { Response } from 'express';
import {
  DomainError,
  NotFoundError,
  ConflictError,
  UnauthorizedError,
  ForbiddenError,
} from '../errors';

@Catch(DomainError)
export class DomainExceptionFilter implements ExceptionFilter {
  catch(exception: DomainError, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const statusCode = this.resolveStatus(exception);

    response.status(statusCode).json({
      statusCode,
      message: exception.message,
    });
  }

  private resolveStatus(exception: DomainError): number {
    if (exception instanceof NotFoundError) return 404;
    if (exception instanceof ConflictError) return 409;
    if (exception instanceof UnauthorizedError) return 401;
    if (exception instanceof ForbiddenError) return 403;
    return 500;
  }
}
