import { Component, OnInit } from '@angular/core';
import { CommonModule, NgOptimizedImage } from '@angular/common';
import { ImageComponent } from '../image/image.component';
import { EnrollResponse, WechatService } from '../../services/wechat.service';
import { isWechat } from '../../core/util';
import { MatButtonModule } from '@angular/material/button';
import { RouterLink } from '@angular/router';
import { AuthService } from '../../core/auth/auth.service';
import { LoginService } from '../../services/login.service';
import { BottomSheetComponent } from '../../shared/bottom-sheet.component';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { EnrollDialogComponent } from '../enroll/enroll-dialog.component';
import { NotificationService } from '../../core/notifications/notification.service';
import { MatCardModule } from '@angular/material/card';

const activityId = 'activity-1'

@Component({
  selector: 'app-article-1',
  standalone: true,
  imports: [CommonModule, NgOptimizedImage, ImageComponent, MatButtonModule, RouterLink, MatDialogModule, MatCardModule],
  templateUrl: './article-1.component.html',
})
export class Article1Component implements OnInit {
  public enroll: EnrollResponse | null = null;

  constructor(
    private wechatService: WechatService,
    public auth: AuthService,
    private loginService: LoginService,
    private dialog: MatDialog,
    private notify: NotificationService,
  ) {
    if (isWechat()) {
      wechatService.refresh(location.href.split('#')[0]).subscribe(() => {
        let title = '常武功夫咏春';
        let desc = '常武功夫极简咏春拳成人班公开课等你来！';
        wechatService.wx.updateAppMessageShareData({
          title: title,
          desc: desc,
          link: location.href,
          imgUrl: 'https://gongfu.grianchan.com/web/assets/article-1/WechatIMG285.jpg', // 分享图标
          success: function () {
            console.log('shared success');
          }
        });
      });
    }
  }

  ngOnInit(): void {
    if (this.auth.isAuthenticated()) {
      this.refreshEnroll()
    }
  }

  // 点击报名按钮
  clickEnroll() {
    // 去登陆
    if (!this.auth.isAuthenticated()) {
      this.dialog.open(BottomSheetComponent, {
        data: {
          true: '登陆',
          false: '取消'
        }
      }).afterClosed().subscribe(v => {
        if (v) {
          this.loginService.setLoginRedirectUrl('/pages/article-1')
          window.location.href = '/wechat-login';
          return
        }
      })
    }

    this.dialog.open(EnrollDialogComponent, {}).afterClosed().subscribe(v => {
      if (v) {
        this.wechatService.pay({
          amount: 1,
          attach: activityId,
          activity_id: activityId,
          description: '常武功夫活动报名',
          phone: v.phone,
          sex: v.sex,
          username: v.username,
          onSuccess: () => {
            this.notify.success('报名成功')
            this.refreshEnroll()
          },
          onFailed: (err: string) => {
            if (err === undefined) {
              this.refreshEnroll()
            } else {
              this.notify.warn('报名失败: ' + err)
            }
          }
        })
      }
    })
  }

  refreshEnroll() {
    this.wechatService.getEnroll(activityId).subscribe(v => {
      if (v.id !== 0) {
        this.enroll = v
      } else {
        this.enroll = null
      }
      this.notify.success('已刷新')
    })
  }

  // 跳转去登陆
  toLogin() {
    this.loginService.setLoginRedirectUrl('/pages/article-1')
    window.location.href = '/wechat-login';
    return false;
  }
}
