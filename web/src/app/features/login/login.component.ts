import {
  ChangeDetectionStrategy,
  Component,
  OnInit,
  ViewChild
} from '@angular/core';
import { MatProgressBar, MatProgressBarModule } from '@angular/material/progress-bar';
import { UntypedFormControl, UntypedFormGroup, Validators, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NotificationService } from '../../core/notifications/notification.service';
import { AuthService } from '../../core/auth/auth.service';
import { ApiService } from '../../api/api.service';
import { Login } from '../../api/models/login';
import { faTrash } from '@fortawesome/free-solid-svg-icons';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { NgOptimizedImage } from '@angular/common';

@Component({
    selector: 'anms-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.scss'],
    changeDetection: ChangeDetectionStrategy.Default,
    standalone: true,
    imports: [MatProgressBarModule, NgOptimizedImage, FormsModule, ReactiveFormsModule, MatFormFieldModule, MatInputModule]
})
export class LoginComponent implements OnInit {
  @ViewChild(MatProgressBar) public progressBar: MatProgressBar | undefined;
  public faTrash = faTrash;
  public sendCodeTime = 0;
  public inSubmit = false;
  public inLoading = false;

  public signinForm: UntypedFormGroup = new UntypedFormGroup({
    phone: new UntypedFormControl('', Validators.required),
    code: new UntypedFormControl('', Validators.required),
    rememberMe: new UntypedFormControl(false)
  });

  constructor(
    private router: Router,
    private auth: AuthService,
    private notify: NotificationService,
    private api: ApiService,
    private activeRoute: ActivatedRoute,
    private authService: AuthService
  ) {}

  ngOnInit() {
    console.log('snapshot', this.activeRoute.snapshot);
    if (!this.authService.isAuthenticated()) {
      const accessToken =
        this.activeRoute.snapshot.queryParamMap.get('accessToken') ?? '';
      if (accessToken !== '') {
        this.authService.login(accessToken);
        this.api.me().subscribe(
          (user) => {
            this.authService.setUser(user);
          },
          () => {},
          () => this.router.navigate(['/me']).then(() => location.reload())
        );
        return;
      }
    } else {
      // 已登录且在登录页，跳转到 /me
      this.toMe();
    }

    setInterval(() => {
      if (this.sendCodeTime > 0) {
        this.sendCodeTime--;
      }
    }, 1000);
  }

  submit() {
    const signinData = this.signinForm.value;
    console.log('submit');
    if (this.progressBar === undefined) {
      return;
    }

    this.progressBar.mode = 'indeterminate';
    this.inLoading = true;
    this.api
      .login(signinData['phone'], signinData['code'] + '')
      .subscribe((data: Login) => {
        if (data.accessToken === undefined) {
          this.notify.error('invalid login response');
          return;
        }
        this.auth.login(data.accessToken);
        this.notify.success('Login success');
        this.toMe();
      })
      .add(() => {
        if (this.progressBar !== undefined) {
          this.progressBar.mode = 'determinate';
        }
        this.inSubmit = false;
        this.inLoading = false;
      });
  }

  // 清空手机号
  cleanPhone() {
    this.signinForm.get('phone')?.setValue('');
  }

  toWechatLogin() {
    window.location.href = '/wechat-login';
    return false;
  }

  toMe() {
    // 这里有个 bug，不刷新页面的话，导航到 /me 后所有链接都点击不正常
    this.router.navigate(['/me']).then(() => location.reload());
  }
}
