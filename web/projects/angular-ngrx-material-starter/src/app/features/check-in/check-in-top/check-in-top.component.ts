import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { CheckIn } from '../../../api/models/check-in';
import { ApiService } from '../../../api/api.service';
import { NotificationService } from '../../../core/notifications/notification.service';
import { Route, Router } from '@angular/router';

@Component({
  selector: 'anms-check-in-top',
  templateUrl: './check-in-top.component.html',
  styleUrls: ['./check-in-top.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class CheckInTopComponent implements OnInit {
  public checkInList: CheckIn[] | undefined

  constructor(
    private api: ApiService,
    private notification: NotificationService,
    private router: Router
  ) {
  }

  ngOnInit(): void {
    const today = new Date()
    let start = 0
    let end = 0
    // 5 点过后才算今天榜单
    if (today.getHours() > 4) {
      // 今天的榜单: 今天凌晨 5:00 - 次日凌晨 4:59
      start = new Date(today.getFullYear(), today.getMonth(), today.getDate(), 5, 0, 0).getTime() / 1000
      today.setDate(today.getDate() + 1)
      const nextDay = today
      end = new Date(nextDay.getFullYear(), nextDay.getMonth(), nextDay.getDate(), 4, 59, 59).getTime() / 1000
    } else {
      // 昨天的榜单: 昨天凌晨 5:00 - 今天凌晨 4:59
      end = new Date(today.getFullYear(), today.getMonth(), today.getDate(), 4, 59, 59).getTime() / 1000
      today.setDate(today.getDate() - 1)
      const beforeDay = today
      start = new Date(beforeDay.getFullYear(), beforeDay.getMonth(), beforeDay.getDate(), 5, 0, 0).getTime() / 1000
    }
    console.log('top time range: ', new Date(start * 1000).toString(), new Date(end * 1000).toString())
    this.api.getCheckInTop(start, end).subscribe(data => this.checkInList = data)
  }
}
