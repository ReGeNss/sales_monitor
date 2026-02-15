import { Controller, Put, Delete, Body, UseGuards } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth } from '@nestjs/swagger';
import { UsersService } from './users.service';
import { RegisteredUserGuard } from '../auth/guards/registered-user.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';
import { User } from '@sales-monitor/database';
import { UpdateNotificationTokenDto } from './dto/update-notification-token.dto';

@ApiTags('users')
@Controller('users')
@UseGuards(RegisteredUserGuard)
@ApiBearerAuth()
export class UsersController {
  constructor(private readonly usersService: UsersService) {}
  
  @Put('notification-token')
  @ApiOperation({ summary: 'Update notification token' })
  @ApiResponse({ status: 200, description: 'Token updated successfully' })
  async updateNotificationToken(
    @CurrentUser() user: User,
    @Body() dto: UpdateNotificationTokenDto,
  ) {
    return this.usersService.updateNotificationToken(user.userId, dto.nfToken);
  }

  @Delete('account')
  @ApiOperation({ summary: 'Delete user account' })
  @ApiResponse({ status: 200, description: 'Account deleted successfully' })
  async deleteAccount(@CurrentUser() user: User) {
    return this.usersService.deleteAccount(user.userId);
  }
}
