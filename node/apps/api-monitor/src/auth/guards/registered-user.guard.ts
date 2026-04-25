import { Injectable, CanActivate, ExecutionContext } from '@nestjs/common';
import { ForbiddenError } from '../../common/errors';

@Injectable()
export class RegisteredUserGuard implements CanActivate {
  canActivate(context: ExecutionContext): boolean {
    const request = context.switchToHttp().getRequest();
    const user = request.user;

    if (!user || user.isGuest) {
      throw new ForbiddenError('This action requires a registered account');
    }

    return true;
  }
}
