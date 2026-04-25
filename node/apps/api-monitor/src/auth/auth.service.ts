import { Injectable } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import * as bcrypt from 'bcrypt';
import { AuthRepository } from './auth.repository';
import { UserDomain } from '../common/domain/user.domain';
import { RegisterDto } from './dto/register.dto';
import { GuestLoginDto } from './dto/guest-login.dto';
import { AuthResponseDto } from './dto/auth-response.dto';

@Injectable()
export class AuthService {
  constructor(
    private readonly authRepository: AuthRepository,
    private readonly jwtService: JwtService,
  ) {}

  async register(registerDto: RegisterDto): Promise<AuthResponseDto> {
    const hashedPassword = await bcrypt.hash(registerDto.password, 10);
    const user = await this.authRepository.createUser(registerDto.login, hashedPassword);
    return this.generateTokenResponse(user);
  }

  async validateUser(login: string, password: string): Promise<UserDomain | null> {
    const result = await this.authRepository.findByLogin(login);
    if (!result) return null;

    const isPasswordValid = await bcrypt.compare(password, result.hashedPassword);
    if (!isPasswordValid) return null;

    return result.user;
  }

  async login(user: UserDomain): Promise<AuthResponseDto> {
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

  private generateTokenResponse(user: UserDomain): AuthResponseDto {
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
