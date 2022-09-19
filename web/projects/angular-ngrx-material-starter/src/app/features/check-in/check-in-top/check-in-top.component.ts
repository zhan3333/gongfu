import { Component, OnInit, ChangeDetectionStrategy, LOCALE_ID } from '@angular/core';
import { CheckIn } from '../../../api/models/check-in';
import { ApiService } from '../../../api/api.service';
import { NotificationService } from '../../../core/notifications/notification.service';
import { ActivatedRoute, Route, Router } from '@angular/router';
import { WechatService } from '../../../services/wechat.service';
import { HttpParams } from '@angular/common/http';
import { formatDate } from '@angular/common';

@Component({
  selector: 'anms-check-in-top',
  templateUrl: './check-in-top.component.html',
  styleUrls: ['./check-in-top.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class CheckInTopComponent implements OnInit {
  public checkInList: CheckIn[] | undefined
  public start = 0
  public end = 0

  constructor(
    private api: ApiService,
    private notification: NotificationService,
    private router: Router,
    private wechatService: WechatService,
    private activeRouter: ActivatedRoute,
  ) {
  }

  ngOnInit(): void {
    const query = this.activeRouter.snapshot.queryParamMap
    if (query.has('startAt') && query.has('endAt')) {
      // 参数中有则使用参数中的数据
      this.start = parseInt(<string>query.get('startAt'), 10)
      this.end = parseInt(<string>query.get('endAt'), 10)
    } else {
      // 参数中没有则计算今天的时间
      const today = new Date()

      // 5 点过后才算今天榜单
      if (today.getHours() > 4) {
        // 今天的榜单: 今天凌晨 5:00 - 次日凌晨 4:59
        this.start = new Date(today.getFullYear(), today.getMonth(), today.getDate(), 5, 0, 0).getTime() / 1000
        today.setDate(today.getDate() + 1)
        const nextDay = today
        this.end = new Date(nextDay.getFullYear(), nextDay.getMonth(), nextDay.getDate(), 4, 59, 59).getTime() / 1000
      } else {
        // 昨天的榜单: 昨天凌晨 5:00 - 今天凌晨 4:59
        this.end = new Date(today.getFullYear(), today.getMonth(), today.getDate(), 4, 59, 59).getTime() / 1000
        today.setDate(today.getDate() - 1)
        const beforeDay = today
        this.start = new Date(beforeDay.getFullYear(), beforeDay.getMonth(), beforeDay.getDate(), 5, 0, 0).getTime() / 1000
      }
    }

    console.log('top time range: ', new Date(this.start * 1000).toString(), new Date(this.end * 1000).toString())
    this.api.getCheckInTop(this.start, this.end).subscribe(data => {
      this.checkInList = data
      this.wechatService.refresh(location.href.split('#')[0]).subscribe(
        () => {
          const link = window.location.origin + `/web/check-in-top?start=${this.start}&end=${this.end}`
          const startDate = formatDate(this.start * 1000, 'Y-M-d H:i:s', 'zh-cn')
          const endDate = formatDate(this.end * 1000, 'Y-M-d H:i:s', 'zh-cn')
          const curDate = formatDate(this.start * 1000, 'Y-M-d', 'zh-cn')
          console.log('shared link: ' + link)
          this.wechatService.wx.updateAppMessageShareData({
            title: '打卡统计',
            desc: `打卡人数: ${this.checkInList?.length}
日期: ${curDate}
时间区间: ${startDate} - ${endDate}`,
            link: link,
            imgUrl: 'https://storage-1313942024.cos.ap-shanghai.myqcloud.com/logo.jpeg', // 分享图标
            success: function () {
              console.log('shared success')
            },
          })
        }
      )
    })
  }
}
