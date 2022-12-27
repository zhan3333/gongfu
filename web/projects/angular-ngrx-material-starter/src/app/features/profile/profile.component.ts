import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AuthService } from '../../core/auth/auth.service';
import { NotificationService } from '../../core/notifications/notification.service';
import { ApiService } from '../../api/api.service';
import { displayLevel } from '../../api/models/user';
import { Profile } from '../../api/models/profile';

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
  ) {
  }

  ngOnInit(): void {
    const uuid = this.activeRoute.snapshot.paramMap.get('uuid')
    if (uuid === '' || uuid === null) {
      this.notificationService.error('invalid uuid')
      return
    }
    this.api.getProfile(uuid).subscribe(data => this.profile = data)
  }
}
