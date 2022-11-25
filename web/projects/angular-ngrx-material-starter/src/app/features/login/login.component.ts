import { ChangeDetectionStrategy, Component, OnInit, ViewChild } from '@angular/core';
import { MatProgressBar } from '@angular/material/progress-bar';
import { MatButton } from '@angular/material/button';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NotificationService } from '../../core/notifications/notification.service';
import { AuthService } from '../../core/auth/auth.service';
import { ApiService } from '../../api/api.service';
import { Login } from '../../api/models/login';

@Component({
  selector: 'anms-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class LoginComponent implements OnInit {
  @ViewChild(MatProgressBar) public progressBar: MatProgressBar | undefined;
  @ViewChild(MatButton) public submitButton: MatButton | undefined;
  public sendCodeTime = 0
  public isBindPhone = false;
  public isLogin = false;

  public signinForm: FormGroup = new FormGroup({
    phone: new FormControl('', Validators.required),
    code: new FormControl('', Validators.required),
    rememberMe: new FormControl(false)
  });

  constructor(
    private router: Router,
    private auth: AuthService,
    private notify: NotificationService,
    private api: ApiService,
    private activeRoute: ActivatedRoute,
    private authService: AuthService,
  ) {
  }

  ngOnInit() {
    console.log('snapshot', this.activeRoute.snapshot)
    if (!this.authService.isAuthenticated()) {
      const accessToken = this.activeRoute.snapshot.queryParamMap.get('accessToken') ?? '';
      if (accessToken !== '') {
        this.authService.login(accessToken)
        this.api.me().subscribe(user => {
            this.authService.setUser(user)
          },
          () => {
          },
          () => this.router.navigate(['/me']).then(() => location.reload())
        )
        return
      }
    }

    if (this.activeRoute.snapshot.queryParamMap.get('type') === 'login') {
      this.isLogin = true
    } else {
      this.isBindPhone = true
    }
    setInterval(() => {
      if (this.sendCodeTime > 0) {
        this.sendCodeTime--
      }
    }, 1000)
  }

  public sendValidCode() {
    const phone = this.signinForm.get('phone')?.value
    if (this.isBindPhone) {
      this.api.sendValidCode(phone).subscribe(
        () => {
          this.notify.success('发送验证码成功')
        },
        error => {
          this.notify.error('发送验证码失败: ' + error)
        }
      )
      this.sendCodeTime = 60
    } else {
      this.notify.warn('暂不支持手机号登录')
    }
  }

  submit() {
    const signinData = this.signinForm.value;
    if (this.submitButton === undefined || this.progressBar === undefined) {
      return
    }

    if (this.isBindPhone) {
      // 绑定
      this.progressBar.mode = 'indeterminate';
      const phone = signinData['phone']
      const code = signinData['code']
      this.api.validCode(phone, code).subscribe(
        () => {
          this.notify.success('绑定成功')
          this.router.navigate(['/me'])
        },
        error => {
          this.notify.error('绑定失败: ' + error)
        },
        () => {
          if (this.submitButton !== undefined && this.progressBar !== undefined) {
            this.submitButton.disabled = false;
            this.progressBar.mode = 'determinate';
          }
        }
      )
    } else {
      // 登录
      this.progressBar.mode = 'indeterminate';
      this.api.login(signinData['phone'], signinData['code'] + '')
        .subscribe(
          (data: Login) => {
            if (data.accessToken === undefined) {
              this.notify.error('invalid login response')
              return
            }
            this.auth.login(data.accessToken)
            this.notify.success('Login success')
            // 这里有个 bug，不刷新页面的话，导航到 /me 后所有链接都点击不正常
            this.router.navigate(['/me']).then(() => location.reload())
          },
        )
        .add(() => {
          if (this.submitButton !== undefined && this.progressBar !== undefined) {
            this.submitButton.disabled = false;
            this.progressBar.mode = 'determinate';
          }
        });
    }
  }

  canSendCode(): boolean {
    if (this.sendCodeTime > 0) {
      return false
    }
    if (!this.signinForm.get('phone')?.valid) {
      return false
    }
    return true
  }

  // 清空手机号
  cleanPhone() {
    this.signinForm.get('phone')?.setValue('')
  }

  toWechatLogin() {
    window.location.href = '/wechat-login'
  }
}
