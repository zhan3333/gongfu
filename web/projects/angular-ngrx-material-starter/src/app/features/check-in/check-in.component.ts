import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { CheckIn } from '../../api/models/check-in';
import { ApiService } from '../../api/api.service';
import { NotificationService } from '../../core/notifications/notification.service';
import { WechatService } from '../../services/wechat.service';

@Component({
  selector: 'anms-check-in',
  templateUrl: './check-in.component.html',
  styleUrls: ['./check-in.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class CheckInComponent implements OnInit {
  public todayCheckIn: CheckIn | undefined
  public file: File | undefined
  public loading = false
  public uploadProgressValue = 0

  constructor(
    private api: ApiService,
    private notification: NotificationService,
    private wechatService: WechatService,
  ) {
  }

  ngOnInit(): void {
    this.getTodayCheckIn()
  }

  public onFileSelected(event: Event) {
    if (event.target == null) {
      return
    }
    const input = event.target as HTMLInputElement
    if (input.files === null || input.files.length === 0) {
      return
    }
    this.file = input.files[0]
  }

  public confirmCheckIn() {
    if (this.file === undefined) {
      this.notification.warn('请选择要上传的文件')
      return
    }
    this.loading = true
    this.api.uploadFile(
      this.file,
      value => this.uploadProgressValue = value,
      key => {
        console.log('upload file return key: ', key)
        this.notification.success('上传文件成功，正在进行打卡')
        // @ts-ignore
        this.api.checkIn(key, this.file.name).subscribe(
          () => {
            this.notification.success('打卡成功')
            this.file = undefined
            this.getTodayCheckIn()
          },
          error => this.notification.error('打卡失败，请稍后重试: ' + error),
          () => {
            this.loading = false
            this.uploadProgressValue = 0
          }
        )
      }
    ).subscribe(
      () => {
      },
      error => {
        this.notification.error('上传文件失败，请稍后重试: ' + error)
        this.loading = false
        this.uploadProgressValue = 0
      },
    )
  }

  public resetCheckIn() {
    this.todayCheckIn = undefined
  }

  private refreshSharedData() {
    if (this.todayCheckIn !== undefined) {
      this.wechatService.refresh(location.href.split('#')[0]).subscribe(
        () => {
          let checkInAtStr = '无'
          if (this.todayCheckIn?.createdAt) {
            checkInAtStr = new Date(this.todayCheckIn?.createdAt * 1000 ?? 0).toLocaleString('chinese', {hour12: false})
          }
          this.wechatService.wx.updateAppMessageShareData({
            title: '今日打卡',
            desc: '打卡时间: ' + checkInAtStr,
            link: window.location.origin + '/web/check-in/' + this.todayCheckIn?.key, // 分享链接，该链接域名或路径必须与当前页面对应的公众号 JS 安全域名一致
            imgUrl: 'https://storage-1313942024.cos.ap-shanghai.myqcloud.com/logo.jpeg', // 分享图标
            success: function () {
              console.log('shared success')
            },
          })
        }
      )
    }
  }

  private getTodayCheckIn() {
    this.api.getTodayCheckIn().subscribe(
      checkInExist => {
        console.log('ret', checkInExist)
        if (checkInExist.exists) {
          this.todayCheckIn = checkInExist.checkIn
          console.log('todayCheckIn', this.todayCheckIn)
          this.refreshSharedData()
        }
      }
    )
  }
}
