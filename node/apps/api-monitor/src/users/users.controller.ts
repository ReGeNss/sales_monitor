import { Controller, Put, Delete, Body, UseGuards } from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
  ApiUnauthorizedResponse,
  ApiForbiddenResponse,
  ApiNotFoundResponse,
  ApiBody,
} from '@nestjs/swagger';
import { UsersService } from './users.service';
import { RegisteredUserGuard } from '../auth/guards/registered-user.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';
import { User } from '@sales-monitor/database';
import { UpdateNotificationTokenDto } from './dto/update-notification-token.dto';
import { MessageResponseDto } from '../common/dto/message-response.dto';

@ApiTags('users')
@Controller('users')
@UseGuards(RegisteredUserGuard)
@ApiBearerAuth()
@ApiUnauthorizedResponse({ description: 'Authorization required (Bearer token)' })
@ApiForbiddenResponse({ description: 'Access restricted to registered users only (not guests)' })
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @Put('notification-token')
  @ApiOperation({
    summary: 'Update notification token',
    description: 'Updates the Firebase/notification token for receiving push notifications. Available only to registered users.',
  })
  @ApiBody({ type: UpdateNotificationTokenDto })
  @ApiResponse({
    status: 200,
    description: 'Notification token updated successfully',
    type: MessageResponseDto,
  })
  @ApiNotFoundResponse({ description: 'User not found' })
  async updateNotificationToken(
    @CurrentUser() user: User,
    @Body() dto: UpdateNotificationTokenDto,
  ) {
    return this.usersService.updateNotificationToken(user.userId, dto.nfToken);
  }

  @Delete('account')
  @ApiOperation({
    summary: 'Delete user account',
    description: 'Permanently deletes the current registered user account. This action is irreversible.',
  })
  @ApiResponse({
    status: 200,
    description: 'Account deleted successfully',
    type: MessageResponseDto,
  })
  @ApiNotFoundResponse({ description: 'User not found' })
  async deleteAccount(@CurrentUser() user: User) {
    return this.usersService.deleteAccount(user.userId);
  }
}
