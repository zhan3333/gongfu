import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../core/auth/auth.service';
import { NotificationService } from '../../core/notifications/notification.service';
import { ApiService } from '../../api/api.service';
import { User } from '../../api/models/user';
import { FormControl, Validators } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'anms-me',
  templateUrl: './me.component.html',
  styleUrls: ['./me.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class MeComponent implements OnInit {
  public accessToken = ''
  public user: User | undefined;
  public bindPhone = new FormControl('', [
    Validators.required,
    Validators.pattern('1(3|4|5|7|8)\\d{9}')
  ])
  public validCode = new FormControl('', [
    Validators.required,
    Validators.maxLength(4),
    Validators.minLength(4)
  ])
  public sendValidCodeLimiting = 0
  public loading = false

  constructor(
    private activeRoute: ActivatedRoute,
    private authService: AuthService,
    private readonly notificationService: NotificationService,
    private api: ApiService,
    private router: Router
  ) {
  }

  ngOnInit(): void {
    if (!this.authService.isAuthenticated()) {
      this.accessToken = this.activeRoute.snapshot.queryParamMap.get('accessToken') ?? '';
      if (this.accessToken === '') {
        this.notificationService.warn('no login')
        return
      } else {
        this.authService.login(this.accessToken)
        this.api.me().subscribe(user => {
          this.authService.setUser(user)
        })
      }
    }
    this.displayUserInfo()
    setInterval(() => {
      if (this.sendValidCodeLimiting > 0) {
        this.sendValidCodeLimiting--
      }
    }, 1000)
  }

  public toBindPhonePage() {
    this.router.navigate(['/login'], {queryParams: {type: 'bind_phone'}})
  }

  public toLoginPage() {
    this.router.navigate(['/login'], {queryParams: {type: 'login'}})
  }

  public toWeChatLogin() {
    window.location.href = '/login'
  }

  private displayUserInfo() {
    this.api.me().subscribe(
      user => this.user = user,
      error => {
        if (error instanceof HttpErrorResponse) {
          if (error.status === 401 || error.status === 403) {
            this.authService.logout()
          }
        }
      })
  }
}
