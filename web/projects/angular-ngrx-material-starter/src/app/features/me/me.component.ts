import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../core/auth/auth.service';
import { NotificationService } from '../../core/notifications/notification.service';
import { ApiService } from '../../api/api.service';
import { ICoach, User } from '../../api/models/user';
import { FormControl, Validators } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import {
  faCalendar,
  faCheck,
  faUser,
  faWrench
} from '@fortawesome/free-solid-svg-icons';
import { faBolt } from '@fortawesome/free-solid-svg-icons/faBolt';

@Component({
  selector: 'anms-me',
  templateUrl: './me.component.html',
  styles: [
    `
      .user_header {
        width: 90px;
        height: 90px;
        border-radius: 20%;
      }
    `
  ],
  changeDetection: ChangeDetectionStrategy.Default
})
export class MeComponent implements OnInit {
  public faUser = faUser;
  public faWrench = faWrench;
  public faCalendar = faCalendar;
  public faBolt = faBolt;
  public faCheck = faCheck;
  public accessToken = '';
  public user: User | undefined;
  public coach: ICoach | undefined;
  public bindPhone = new FormControl('', [
    Validators.required,
    Validators.pattern('1(3|4|5|7|8)\\d{9}')
  ]);
  public validCode = new FormControl('', [
    Validators.required,
    Validators.maxLength(4),
    Validators.minLength(4)
  ]);
  public sendValidCodeLimiting = 0;
  public loading = false;

  constructor(
    private activeRoute: ActivatedRoute,
    private authService: AuthService,
    private readonly notificationService: NotificationService,
    private api: ApiService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.displayUserInfo();
    setInterval(() => {
      if (this.sendValidCodeLimiting > 0) {
        this.sendValidCodeLimiting--;
      }
    }, 1000);
  }

  public toBindPhonePage() {
    this.router.navigate(['/login'], { queryParams: { type: 'bind_phone' } });
  }

  public toLoginPage() {
    this.router.navigate(['/login'], { queryParams: { type: 'login' } });
  }

  public toWeChatLogin() {
    window.location.href = '/login';
  }

  private displayUserInfo() {
    this.api.me().subscribe(
      (user) => {
        this.user = user;
        this.authService.setUser(user);
        if (user.hasRole('coach')) {
          this.api.getCoach().subscribe((data) => (this.coach = data));
        }
      },
      (error) => {
        if (error instanceof HttpErrorResponse) {
          if (error.status === 401 || error.status === 403) {
            this.authService.logout();
            return;
          }
        }
      }
    );
  }
}
