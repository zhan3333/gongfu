import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { NotificationService } from '../../../core/notifications/notification.service';
import { ApiService, CheckInComment } from '../../../api/api.service';
import { CheckIn } from '../../../api/models/check-in';
import { refreshSharedCheckInToWechat } from '../../../core/util';
import { WechatService } from '../../../services/wechat.service';
import { DatePipe, NgForOf, NgIf, NgOptimizedImage, NgStyle } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { ROLE_ADMIN, ROLE_COACH, User } from '../../../api/models/user';
import { AuthService } from '../../../core/auth/auth.service';
import { FormControl, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatRippleModule } from '@angular/material/core';

@Component({
  selector: 'anms-check-in-show',
  templateUrl: './check-in-show.component.html',
  styleUrls: ['./check-in-show.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default,
  standalone: true,
  imports: [MatCardModule, NgIf, RouterLink, DatePipe, NgForOf, NgOptimizedImage, ReactiveFormsModule, MatInputModule, MatButtonModule, NgStyle, MatRippleModule]
})
export class CheckInShowComponent implements OnInit {
  public checkIn: CheckIn | undefined
  comments: CheckInComment[] = [];
  user: User
  commentForm = new FormControl('', {nonNullable: true, validators: [Validators.required]});
  protected readonly ROLE_COACH = ROLE_COACH;
  protected readonly ROLE_ADMIN = ROLE_ADMIN;

  constructor(
    private activeRoute: ActivatedRoute,
    private notification: NotificationService,
    private api: ApiService,
    private wechatService: WechatService,
    private auth: AuthService
  ) {
    this.user = auth.getUser()
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
        this.refreshComments()
      })
    }
  }

  submitComment() {
    if (this.commentForm.invalid) {
      return
    }
    this.api.createCheckInComment(this.checkIn!.id!, this.commentForm.value).subscribe({
      next: () => {
        this.commentForm.reset()
        this.refreshComments()
        this.notification.success('ok')
      }
    })
  }

  private refreshComments() {
    this.api.getCheckInComments(this.checkIn!.id!).subscribe(data => {
      this.comments = data
    })
  }
}
