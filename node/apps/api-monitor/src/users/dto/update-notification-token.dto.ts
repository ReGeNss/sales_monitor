import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsOptional } from 'class-validator';

export class UpdateNotificationTokenDto {
  @ApiProperty({ description: 'Firebase/notification token', required: false })
  @IsString()
  @IsOptional()
  nfToken?: string;
}
