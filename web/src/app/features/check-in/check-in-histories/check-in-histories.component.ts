import { Component, OnInit, ChangeDetectionStrategy, ViewEncapsulation, ViewChild } from '@angular/core';
import { ApiService } from '../../../api/api.service';
import { ActivatedRoute } from '@angular/router';
import { NotificationService } from '../../../core/notifications/notification.service';
import { MatCalendar, MatCalendarCellCssClasses, MatDatepickerModule } from '@angular/material/datepicker';
import { CheckInList } from '../../../api/models/check-in';
import * as moment from 'moment';
import { CalendarHeaderComponent } from '../../../shared/calendar-header/calendar-header.component';
import { MatCardModule } from '@angular/material/card';

@Component({
    selector: 'anms-check-in-histories',
    templateUrl: './check-in-histories.component.html',
    styleUrls: ['./check-in-histories.component.scss'],
    changeDetection: ChangeDetectionStrategy.Default,
    encapsulation: ViewEncapsulation.None,
    standalone: true,
    imports: [MatCardModule, MatDatepickerModule],
})
export class CheckInHistoriesComponent implements OnInit {
  @ViewChild(MatCalendar) calendar: MatCalendar<Date> | undefined;
  public calendarHeader = CalendarHeaderComponent

  public userID = 0
  public checkInList: CheckInList = []

  constructor(
    private api: ApiService,
    private route: ActivatedRoute,
    private notification: NotificationService
  ) {
  }

  ngOnInit(): void {
    this.userID = parseInt(<string>this.route.snapshot.queryParamMap.get('userID'), 10)
    console.log('get user id', this.userID, this.route.snapshot.queryParamMap.get('userID'))
    if (this.userID === 0) {
      this.notification.error('无效的 userID 参数')
      return
    }

    this.onMonthSelected(moment().toDate())
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

  public onMonthSelected(date: Date | null) {
    console.log('select month', date)
    const startDate = moment(date).startOf('month').format('YYYYMMDD')
    const endDate = moment(date).endOf('month').format('YYYYMMDD')
    console.log('dates', startDate, endDate)
    this.api.getCheckInHistories(this.userID, startDate, endDate).subscribe((data: CheckInList) => {
      console.log('get histories', data)
      this.checkInList = data
      if (this.calendar !== undefined) {
        this.calendar.updateTodaysDate()
      }
      this.notification.success('加载成功')
    })
  }
}
