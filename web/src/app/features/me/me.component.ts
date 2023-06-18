import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnInit
} from '@angular/core';
import { ActivatedRoute, Params, Router, RouterLink } from '@angular/router';
import { AuthService } from '../../core/auth/auth.service';
import { NotificationService } from '../../core/notifications/notification.service';
import { ApiService } from '../../api/api.service';
import { ICoach, User } from '../../api/models/user';
import { UntypedFormControl, Validators } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import {
  faAngleRight,
  faCalendar,
  faCheck,
  faUser,
  faWrench,
  IconDefinition
} from '@fortawesome/free-solid-svg-icons';
import { faBolt } from '@fortawesome/free-solid-svg-icons/faBolt';
import { NgForOf, NgIf, NgOptimizedImage, NgStyle } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { MatChipsModule } from '@angular/material/chips';
import { MatRippleModule } from '@angular/material/core';

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
    NgStyle,
    MatChipsModule,
    MatRippleModule,
    NgForOf
  ],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MeComponent implements OnInit {
  public accessToken = '';
  public user: User | undefined = undefined;
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
  protected readonly faAngleRight = faAngleRight;

  constructor(
    private activeRoute: ActivatedRoute,
    private authService: AuthService,
    private readonly notificationService: NotificationService,
    private api: ApiService,
    private router: Router,
    private changeRef: ChangeDetectorRef
  ) {}

  public blockMenus(user: User): {
    img: string;
    name: string;
    link: string;
  }[] {
    return [
      {
        img: 'https://lanhu.oss-cn-beijing.aliyuncs.com/pspcvyjzpfl68le8lxb1vmx8a0q8q3tlht8965c932-bc41-4dcc-a256-f1eefd950ced',
        name: '去打卡',
        link: '/check-in'
      },
      {
        img: 'https://lanhu.oss-cn-beijing.aliyuncs.com/pse062snx0cmqaf53pe9abguf9wzv373lse743c141-d793-4024-9e58-57702bd4a8d6',
        name: '打卡排行',
        link: '/check-in-top'
      },
      {
        img: 'https://lanhu.oss-cn-beijing.aliyuncs.com/ps1ve9zudnkd8t0f63uxinsolqq479qy5j8e06716b-4dd0-491d-bfe6-bce982c080f8',
        name: '历史打卡',
        link: '/check-in-histories'
      },
      {
        img: 'https://lanhu.oss-cn-beijing.aliyuncs.com/psqrytvmyg29m8x1bajsr5zkl1jxycqkzun7f81ffc5-78bb-4b26-ab59-0ff4f8228cd8',
        name: '更多信息',
        link: `/profile/${user.uuid}`
      }
    ];
  }

  public listMenus(user: User): {
    icon: IconDefinition;
    name: string;
    link: string;
    queryParams: Params | null;
  }[] {
    return [
      {
        icon: faCheck,
        name: '去打卡',
        link: '/check-in',
        queryParams: null
      },
      {
        icon: faBolt,
        name: '我的历史打卡',
        link: '/check-in-histories',
        queryParams: { userID: user.id }
      },
      {
        icon: faCalendar,
        name: '打卡榜单',
        link: '/check-in-top',
        queryParams: null
      },
      {
        icon: faUser,
        name: '更多信息',
        link: `/profile/${user.uuid}`,
        queryParams: null
      },
      {
        icon: faWrench,
        name: '设置',
        link: '/setting',
        queryParams: null
      }
    ];
  }

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
        console.log('user', user);
        this.user = user;
        this.authService.setUser(user);
        if (user.hasRole('coach')) {
          this.api.getCoach().subscribe((data) => (this.coach = data));
        }
        this.changeRef.markForCheck();
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
