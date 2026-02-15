import { Injectable, ConflictException, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { EntityManager } from '@mikro-orm/core';
import { User } from '@sales-monitor/database';
import * as bcrypt from 'bcrypt';
import { RegisterDto } from './dto/register.dto';
import { GuestLoginDto } from './dto/guest-login.dto';
import { AuthResponseDto } from './dto/auth-response.dto';

@Injectable()
export class AuthService {
  constructor(
    private readonly em: EntityManager,
    private readonly jwtService: JwtService,
  ) {}

  async register(registerDto: RegisterDto): Promise<AuthResponseDto> {
    const existingUser = await this.em.findOne(User, { login: registerDto.login });
    if (existingUser) {
      throw new ConflictException('User with this login already exists');
    }

    const hashedPassword = await bcrypt.hash(registerDto.password, 10);
    const user = new User();
    user.login = registerDto.login;
    user.password = hashedPassword;

    await this.em.persistAndFlush(user);

    return this.generateTokenResponse(user);
  }

  async validateUser(login: string, password: string): Promise<User | null> {
    const user = await this.em.findOne(User, { login });
    if (!user) {
      return null;
    }

    const isPasswordValid = await bcrypt.compare(password, user.password);
    if (!isPasswordValid) {
      return null;
    }

    return user;
  }

  async login(user: User): Promise<AuthResponseDto> {
    return this.generateTokenResponse(user);
  }

  guestLogin(dto: GuestLoginDto): { access_token: string } {
    const payload = {
      isGuest: true,
      deviceModel: dto.deviceModel,
      appVersion: dto.appVersion,
      platform: dto.platform,
      locale: dto.locale,
    };
    return { access_token: this.jwtService.sign(payload) };
  }

  private generateTokenResponse(user: User): AuthResponseDto {
    const payload = { sub: user.userId, login: user.login };
    return {
      access_token: this.jwtService.sign(payload),
      user: {
        userId: user.userId,
        login: user.login,
      },
    };
  }
}
