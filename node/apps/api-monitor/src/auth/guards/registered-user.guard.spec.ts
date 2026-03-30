import { ForbiddenException, ExecutionContext } from '@nestjs/common';
import { RegisteredUserGuard } from './registered-user.guard';

const buildContext = (user: any): ExecutionContext =>
  ({
    switchToHttp: () => ({
      getRequest: () => ({ user }),
    }),
  }) as any;

describe('RegisteredUserGuard', () => {
  let guard: RegisteredUserGuard;

  beforeEach(() => {
    guard = new RegisteredUserGuard();
  });

  it('allows a registered user (isGuest: false)', () => {
    const ctx = buildContext({ userId: 1, login: 'alice', isGuest: false });
    expect(guard.canActivate(ctx)).toBe(true);
  });

  it('allows a registered user without isGuest property', () => {
    const ctx = buildContext({ userId: 1, login: 'alice' });
    expect(guard.canActivate(ctx)).toBe(true);
  });

  it('throws ForbiddenException for a guest user (isGuest: true)', () => {
    const ctx = buildContext({ isGuest: true, deviceModel: 'iPhone' });
    expect(() => guard.canActivate(ctx)).toThrow(ForbiddenException);
  });

  it('throws ForbiddenException when user is undefined', () => {
    const ctx = buildContext(undefined);
    expect(() => guard.canActivate(ctx)).toThrow(ForbiddenException);
  });

  it('throws ForbiddenException when user is null', () => {
    const ctx = buildContext(null);
    expect(() => guard.canActivate(ctx)).toThrow(ForbiddenException);
  });

  it('includes the correct message in ForbiddenException', () => {
    const ctx = buildContext({ isGuest: true });
    expect(() => guard.canActivate(ctx)).toThrow('This action requires a registered account');
  });
});
