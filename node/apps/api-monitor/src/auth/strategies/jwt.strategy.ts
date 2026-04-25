import { Injectable } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';
import { UsersRepository } from '../../users/users.repository';
import { UnauthorizedError } from '../../common/errors';

export interface JwtPayload {
  sub?: number;
  login?: string;
  isGuest?: boolean;
  deviceModel?: string;
  appVersion?: string;
  platform?: string;
  locale?: string;
}

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  constructor(private readonly usersRepository: UsersRepository) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: false,
      secretOrKey: process.env.JWT_SECRET || 'your-secret-key-change-in-production',
    });
  }

  async validate(payload: JwtPayload) {
    if (payload.isGuest) {
      return {
        isGuest: true,
        deviceModel: payload.deviceModel,
        appVersion: payload.appVersion,
        platform: payload.platform,
        locale: payload.locale,
      };
    }

    try {
      return await this.usersRepository.findById(payload.sub!);
    } catch {
      throw new UnauthorizedError('User not found');
    }
  }
}
