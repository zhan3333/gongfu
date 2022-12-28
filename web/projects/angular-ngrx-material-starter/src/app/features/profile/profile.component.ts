import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { NotificationService } from '../../core/notifications/notification.service';
import { ApiService } from '../../api/api.service';
import { displayLevel } from '../../api/models/user';
import { Profile } from '../../api/models/profile';
import { refreshSharedProfileToWechat } from '../../core/util';
import { WechatService } from '../../services/wechat.service';

@Component({
  selector: 'anms-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class ProfileComponent implements OnInit {
  public profile: Profile | undefined;
  public displayLevel = displayLevel

  constructor(
    private activeRoute: ActivatedRoute,
    private notificationService: NotificationService,
    private api: ApiService,
    private wechat: WechatService,
  ) {
  }

  ngOnInit(): void {
    const uuid = this.activeRoute.snapshot.paramMap.get('uuid')
    if (uuid === '' || uuid === null) {
      this.notificationService.error('invalid uuid')
      return
    }
    this.api.getProfile(uuid).subscribe(data => {
      this.profile = data
      refreshSharedProfileToWechat(this.wechat, this.profile)
    })
  }
}
