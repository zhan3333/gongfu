import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { CheckIn, CheckInCountList } from '../../../api/models/check-in';
import { ApiService } from '../../../api/api.service';
import { NotificationService } from '../../../core/notifications/notification.service';
import { ActivatedRoute, Router } from '@angular/router';
import { WechatService } from '../../../services/wechat.service';
import * as moment from 'moment';
import { FormControl } from '@angular/forms';
import { isWechat } from '../../../core/util';
import { SharedModule } from '../../../shared/shared.module';

@Component({
  selector: 'anms-check-in-top',
  templateUrl: './check-in-top.component.html',
  styleUrls: ['./check-in-top.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class CheckInTopComponent implements OnInit {
  public checkInList: CheckIn[] | undefined
  public checkInContinuousList: CheckInCountList | undefined;
  public checkInCountList: CheckInCountList | undefined;
  public loading = false
  public selectDate = new FormControl('')

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
    let date = ''
    if (query.has('date')) {
      date = <string>query.get('date')
    } else {
      date = this.getTodayDate()
    }
    this.selectDate.valueChanges.subscribe(
      d => this.refresh(moment(d).format('YYYYMMDD'))
    )
    this.selectDate.setValue(moment(date, 'YYYYMMDD'))
    console.log('query date=' + date)

    this.api.getCheckInCountTop().subscribe(data => this.checkInCountList = data)
    this.api.getCheckInContinuousTop().subscribe(data => this.checkInContinuousList = data)
  }

  public get date() {
    return moment(this.selectDate.value as string).format('YYYYMMDD')
  }

  public clickToday() {
    if (!this.isToday()) {
      this.selectDate.setValue(moment(this.getTodayDate(), 'YYYYMMDD'))
    }
  }

  public isToday(): boolean {
    return this.getTodayDate() === moment(this.selectDate.value as string).format('YYYYMMDD')
  }

  private refresh(date: string) {
    this.loading = true
    this.api.getCheckInTop(date).subscribe(
      data => {
        this.checkInList = data
        const shareData = {
          title: '打卡统计',
          desc: `打卡人数: ${this.checkInList?.length || 0}
日期: ${this.date}`,
          link: window.location.origin + `/web/check-in-top?date=${this.date}`,
          imgUrl: 'https://storage-1313942024.cos.ap-shanghai.myqcloud.com/logo.jpeg', // 分享图标
          success: function () {
            console.log('shared success')
          },
        }
        console.log('shareData', shareData)
        if (isWechat()) {
          this.wechatService
            .refresh(location.href.split('#')[0])
            .subscribe(() => this.wechatService.wx.updateAppMessageShareData(shareData))
        }
      },
      () => {
      },
      () => this.loading = false
    )
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
