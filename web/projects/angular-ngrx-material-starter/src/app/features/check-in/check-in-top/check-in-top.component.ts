import { Component, OnInit, ChangeDetectionStrategy, LOCALE_ID } from '@angular/core';
import { CheckIn } from '../../../api/models/check-in';
import { ApiService } from '../../../api/api.service';
import { NotificationService } from '../../../core/notifications/notification.service';
import { ActivatedRoute, Route, Router } from '@angular/router';
import { WechatService } from '../../../services/wechat.service';
import { HttpParams } from '@angular/common/http';
import { formatDate } from '@angular/common';
import * as moment from 'moment';

@Component({
  selector: 'anms-check-in-top',
  templateUrl: './check-in-top.component.html',
  styleUrls: ['./check-in-top.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class CheckInTopComponent implements OnInit {
  public checkInList: CheckIn[] | undefined
  public date = ''

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
    if (query.has('date')) {
      this.date = <string>query.get('date')
    } else {
      this.date = this.getTodayDate()
    }
    console.log('query date=' + this.date)
    this.refresh()
  }

  public clickToday() {
    if (!this.isToday()) {
      this.date = this.getTodayDate()
      this.refresh()
    }
  }

  public isToday(): boolean {
    return this.getTodayDate() === this.date
  }

  private refresh() {
    this.api.getCheckInTop(this.date).subscribe(data => {
      this.checkInList = data
      this.wechatService.refresh(location.href.split('#')[0]).subscribe(
        () => {
          const link = window.location.origin + `/web/check-in-top?date=${this.date}`
          console.log('shared link: ' + link)
          this.wechatService.wx.updateAppMessageShareData({
            title: '打卡统计',
            desc: `打卡人数: ${this.checkInList?.length}
日期: ${this.date}`,
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

  private getTodayDate(): string {
    const now = moment()
    if (now.hours() >= 5) {
      return now.format('YYYYMMDD')
    } else {
      now.add({days: -1})
      return now.format('YYYYMMDD')
    }
  }
}
