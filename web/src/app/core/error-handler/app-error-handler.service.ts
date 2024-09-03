import { ErrorHandler, Injectable } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';

import { NotificationService } from '../notifications/notification.service';

/** Application-wide error handler that adds a UI notification to the error handling
 * provided by the default Angular ErrorHandler.
 */
@Injectable()
export class AppErrorHandler extends ErrorHandler {
  constructor(private notificationsService: NotificationService) {
    super();
  }

  override handleError(error: Error | HttpErrorResponse) {
    let displayMessage = '';
    if (error instanceof HttpErrorResponse) {
      switch (error.status) {
        case 400:
          if (error.error !== null && error.error['msg']) {
            displayMessage += error.error['msg'];
          } else {
            displayMessage += `客户端错误: ${error.message}`;
          }
          break;
        case 500:
          if (error.error !== null && error.error['msg']) {
            displayMessage += '服务器错误: ';
            displayMessage += error.error['msg'];
          } else {
            displayMessage += `服务器错误: ${error.message}`;
          }
          break;
        default:
          displayMessage += `未知的网络请求错误: ${error.message}`;
          break;
      }
    } else {
      // 其他类型错误
      displayMessage = `有错误发生: ${error.name} / ${error.message}`;
      // 保存到 localstorage 中
      localStorage.setItem("error", JSON.stringify({
        name: error.name,
        stack: error.stack,
        cause: error.cause,
        message: error.message,
      }));
    }

    this.notificationsService.error(displayMessage);

    super.handleError(error);
  }
}
