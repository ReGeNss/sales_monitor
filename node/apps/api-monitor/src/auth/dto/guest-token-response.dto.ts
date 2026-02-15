import { ApiProperty } from '@nestjs/swagger';

export class GuestTokenResponseDto {
  @ApiProperty({ example: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...', description: 'JWT token for guest access' })
  access_token!: string;
}
