import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { AuthService } from '../../core/auth/auth.service';
import { NotificationService } from '../../core/notifications/notification.service';
import { ApiService } from '../../api/api.service';
import { ICoach, User } from '../../api/models/user';
import { UntypedFormControl, Validators } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import {
  faCalendar,
  faCheck,
  faUser,
  faWrench
} from '@fortawesome/free-solid-svg-icons';
import { faBolt } from '@fortawesome/free-solid-svg-icons/faBolt';
import { NgIf, NgOptimizedImage, NgStyle } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

@Component({
  standalone: true,
  selector: 'anms-me',
  templateUrl: './me.component.html',
  imports: [
    MatCardModule,
    MatButtonModule,
    NgOptimizedImage,
    RouterLink,
    MatListModule,
    FontAwesomeModule,
    MatIconModule,
    NgIf,
    NgStyle
  ],
  changeDetection: ChangeDetectionStrategy.Default
})
export class MeComponent implements OnInit {
  public accessToken = '';
  public user: User | undefined;
  public coach: ICoach | undefined;
  public bindPhone = new UntypedFormControl('', [
    Validators.required,
    Validators.pattern('1(3|4|5|7|8)\\d{9}')
  ]);
  public validCode = new UntypedFormControl('', [
    Validators.required,
    Validators.maxLength(4),
    Validators.minLength(4)
  ]);
  public sendValidCodeLimiting = 0;
  public loading = false;
  protected readonly faWrench = faWrench;
  protected readonly faCalendar = faCalendar;
  protected readonly faUser = faUser;
  protected readonly faBolt = faBolt;
  protected readonly faCheck = faCheck;

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
