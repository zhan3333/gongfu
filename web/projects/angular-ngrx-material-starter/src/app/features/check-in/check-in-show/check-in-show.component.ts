import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { NotificationService } from '../../../core/notifications/notification.service';
import { ApiService } from '../../../api/api.service';
import { CheckIn } from '../../../api/models/check-in';

@Component({
  selector: 'anms-check-in-show',
  templateUrl: './check-in-show.component.html',
  styleUrls: ['./check-in-show.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class CheckInShowComponent implements OnInit {
  public checkIn: CheckIn | undefined

  constructor(
    private activeRoute: ActivatedRoute,
    private notification: NotificationService,
    private api: ApiService,
  ) {

  }

  ngOnInit(): void {
    const key = this.activeRoute.snapshot.paramMap.get('key')
    if (key === null) {
      this.notification.error('未获取到 key 参数')
      return
    } else {
      this.api.getCheckInByKey(key).subscribe(data => this.checkIn = data)
    }
  }
}
