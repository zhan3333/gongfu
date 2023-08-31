import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { NotificationService } from '../../../core/notifications/notification.service';
import { ApiService } from '../../../api/api.service';
import { CheckIn } from '../../../api/models/check-in';
import { isWechat, refreshSharedCheckInToWechat } from '../../../core/util';
import { WechatService } from '../../../services/wechat.service';
import { NgIf, DatePipe } from '@angular/common';
import { MatCardModule } from '@angular/material/card';

@Component({
    selector: 'anms-check-in-show',
    templateUrl: './check-in-show.component.html',
    styleUrls: ['./check-in-show.component.scss'],
    changeDetection: ChangeDetectionStrategy.Default,
    standalone: true,
    imports: [MatCardModule, NgIf, RouterLink, DatePipe]
})
export class CheckInShowComponent implements OnInit {
  public checkIn: CheckIn | undefined

  constructor(
    private activeRoute: ActivatedRoute,
    private notification: NotificationService,
    private api: ApiService,
    private wechatService: WechatService,
  ) {

  }

  ngOnInit(): void {
    const key = this.activeRoute.snapshot.paramMap.get('key')
    if (key === null) {
      this.notification.error('未获取到 key 参数')
      return
    } else {
      this.api.getCheckInByKey(key).subscribe(data => {
        this.checkIn = data
        refreshSharedCheckInToWechat(this.wechatService, this.checkIn)
      })
    }
  }
}
