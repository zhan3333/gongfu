import { ChangeDetectionStrategy, Component, OnInit, ViewChild, ViewEncapsulation } from '@angular/core';
import { CheckIn, CheckInList } from '../../api/models/check-in';
import { ApiService } from '../../api/api.service';
import { NotificationService } from '../../core/notifications/notification.service';
import { WechatService } from '../../services/wechat.service';
import { MatCalendar, MatCalendarCellCssClasses } from '@angular/material/datepicker';
import * as moment from 'moment/moment';
import { AuthService } from '../../core/auth/auth.service';
import { CalendarHeaderComponent } from '../../shared/calendar-header/calendar-header.component';
import { isWechat } from '../../core/util';
import { Router } from '@angular/router';

@Component({
  selector: 'anms-check-in',
  templateUrl: './check-in.component.html',
  styleUrls: ['./check-in.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default,
  encapsulation: ViewEncapsulation.None,
})
export class CheckInComponent implements OnInit {
  @ViewChild(MatCalendar) calendar: MatCalendar<Date> | undefined;
  public calendarHeader = CalendarHeaderComponent
  public todayCheckIn: CheckIn | undefined
  public file: File | undefined
  public loading = false
  public uploadProgressValue = 0

  public userID = 0
  public checkInList: CheckInList = []
  public calendarLoading = false

  constructor(
    private api: ApiService,
    private notification: NotificationService,
    private wechatService: WechatService,
    private auth: AuthService,
    private router: Router
  ) {
  }

  ngOnInit(): void {
    this.getTodayCheckIn()
    this.auth.user$.subscribe(user => {
      if (user.id) {
        this.userID = user.id
        this.onMonthSelected(moment().toDate())
      }
    })
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
          error => this.notification.error('打卡失败，请稍后重试: ' + JSON.stringify(error)),
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
        this.notification.error('上传文件失败，请稍后重试: ' + JSON.stringify(error))
        this.loading = false
        this.uploadProgressValue = 0
      },
    )
  }

  public resetCheckIn() {
    this.todayCheckIn = undefined
  }

  public onMonthSelected(date: Date | null) {
    this.calendarLoading = true
    console.log('select month', date)
    const startDate = moment(date).startOf('month').format('YYYYMMDD')
    const endDate = moment(date).endOf('month').format('YYYYMMDD')
    console.log('dates', startDate, endDate)
    this.api.getCheckInHistories(this.userID, startDate, endDate).subscribe(
      (data: CheckInList) => {
        console.log('get histories', data)
        this.checkInList = data
        if (this.calendar !== undefined) {
          this.calendar.updateTodaysDate()
        }
      },
      () => {
      },
      () => this.calendarLoading = false
    )
  }

  public dateClass() {
    return (date: Date): MatCalendarCellCssClasses => {
      const ds = moment(date).format('YYYYMMDD')
      for (const checkIn of this.checkInList) {
        if (checkIn.date === ds) {
          return 'special-date'
        }
      }
      return '';
    }
  }


  public calendarClick() {
    console.log('calendarClick')
  }

  toShowCheckIn(key: string | undefined) {
    if (key === undefined) {
      return
    }
    this.router.navigate(['/check-in', key])
  }

  private refreshSharedData() {
    if (this.todayCheckIn !== undefined && isWechat()) {
      this.wechatService.refresh(location.href.split('#')[0]).subscribe(
        () => {
          let checkInAtStr = '无'
          if (this.todayCheckIn?.createdAt) {
            checkInAtStr = new Date(this.todayCheckIn?.createdAt * 1000 ?? 0).toLocaleString('chinese', {hour12: false})
          }
          this.wechatService.wx.updateAppMessageShareData({
            title: this.todayCheckIn?.userName + ' 的今日打卡',
            desc: '打卡时间: ' + checkInAtStr,
            link: window.location.origin + '/web/check-in/' + this.todayCheckIn?.key, // 分享链接，该链接域名或路径必须与当前页面对应的公众号 JS 安全域名一致
            imgUrl: this.todayCheckIn?.headImgUrl, // 分享图标
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
        if (checkInExist.exists) {
          this.todayCheckIn = checkInExist.checkIn

          this.refreshSharedData()
          this.onMonthSelected(moment().toDate())
        }
      }
    )
  }


}
