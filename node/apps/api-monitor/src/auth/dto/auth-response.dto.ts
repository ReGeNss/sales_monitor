import { ApiProperty } from '@nestjs/swagger';

export class UserResponseDto {
  @ApiProperty()
  userId!: number;

  @ApiProperty()
  login!: string;
}

export class AuthResponseDto {
  @ApiProperty()
  access_token!: string;

  @ApiProperty({ type: UserResponseDto })
  user!: UserResponseDto;
}
