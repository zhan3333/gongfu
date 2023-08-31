import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { NotificationService } from '../../core/notifications/notification.service';
import { ApiService } from '../../api/api.service';
import { Profile } from '../../api/models/profile';
import { refreshSharedProfileToWechat } from '../../core/util';
import { WechatService } from '../../services/wechat.service';
import { displayLevel } from '../../services/coach-level';
import { MatLineModule } from '@angular/material/core';
import { NgIf, NgOptimizedImage, NgFor } from '@angular/common';

@Component({
    selector: 'anms-profile',
    templateUrl: './profile.component.html',
    styleUrls: ['./profile.component.scss'],
    changeDetection: ChangeDetectionStrategy.Default,
    standalone: true,
    imports: [NgIf, NgOptimizedImage, NgFor, MatLineModule]
})
export class ProfileComponent implements OnInit {
  public profile: Profile | undefined;
  public displayLevel = displayLevel;

  constructor(
    private activeRoute: ActivatedRoute,
    private notificationService: NotificationService,
    private api: ApiService,
    private wechat: WechatService
  ) {}

  ngOnInit(): void {
    this.activeRoute.paramMap.subscribe((data) => {
      const uuid = data.get('uuid');
      if (uuid === '' || uuid === null) {
        this.notificationService.error('invalid uuid');
        return;
      }
      this.api.getProfile(uuid).subscribe((profile) => {
        this.profile = profile;
        refreshSharedProfileToWechat(this.wechat, this.profile);
      });
    });
  }
}
