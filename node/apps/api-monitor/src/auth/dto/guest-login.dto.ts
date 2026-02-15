import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsNotEmpty, MaxLength } from 'class-validator';

export class GuestLoginDto {
  @ApiProperty({ example: 'iPhone 15 Pro', description: 'Device model' })
  @IsString()
  @IsNotEmpty()
  @MaxLength(255)
  deviceModel!: string;

  @ApiProperty({ example: '1.0.0', description: 'Application version' })
  @IsString()
  @IsNotEmpty()
  @MaxLength(50)
  appVersion!: string;

  @ApiProperty({ example: 'ios', description: 'Platform (ios, android, web)' })
  @IsString()
  @IsNotEmpty()
  @MaxLength(50)
  platform!: string;

  @ApiProperty({ example: 'uk-UA', description: 'User locale' })
  @IsString()
  @IsNotEmpty()
  @MaxLength(20)
  locale!: string;
}
