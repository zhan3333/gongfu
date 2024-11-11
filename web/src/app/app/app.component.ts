import browser from 'browser-detect';
import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { select, Store } from '@ngrx/store';
import { Observable } from 'rxjs';

import { environment as env } from '../../environments/environment';

import {
  AppState,
  LocalStorageService,
  routeAnimations,
  selectEffectiveTheme,
  selectSettingsLanguage,
  selectSettingsStickyHeader
} from '../core/core.module';
import {
  actionSettingsChangeAnimationsPageDisabled,
  actionSettingsChangeLanguage
} from '../core/settings/settings.actions';
import { AuthService } from '../core/auth/auth.service';
import { ActivatedRoute, NavigationEnd, Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { CHECK_IN_TOP_PATH } from '../core/router/route-path';
import { UserClass } from '../api/models/userClass';
import {
  faBars,
  faCog,
  faPlayCircle,
  faPowerOff,
  faRocket,
  faUser,
  faUserCircle,
  faWrench
} from '@fortawesome/free-solid-svg-icons';
import { TranslateModule } from '@ngx-translate/core';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { MatButtonModule } from '@angular/material/button';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatListModule } from '@angular/material/list';
import { MatSidenavModule } from '@angular/material/sidenav';
import { AsyncPipe, NgFor, NgIf } from '@angular/common';

interface NavigationItem {
  label: string
  link: string
  roles?: string[]
}

@Component({
  selector: 'anms-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  animations: [routeAnimations],
  changeDetection: ChangeDetectionStrategy.Default,
  standalone: true,
  imports: [NgIf, MatSidenavModule, MatListModule, NgFor, RouterLinkActive, RouterLink, MatToolbarModule, MatButtonModule, FontAwesomeModule, MatMenuModule, MatIconModule, RouterOutlet, MatTooltipModule, AsyncPipe, TranslateModule]
})
export class AppComponent implements OnInit {
  faBars = faBars;
  faUserCircle = faUserCircle;
  faPowerOff = faPowerOff;
  faCog = faCog;
  faWrench = faWrench;
  faUser = faUser;
  faPlayCircle = faPlayCircle;
  faRocket = faRocket;
  authLayout = false;
  displayFooter = false;
  isProd = env.production;
  envName = env.envName;
  version = env.versions.app;
  year = new Date().getFullYear();
  logo = 'assets/logo.png';
  languages = ['en', 'de', 'sk', 'fr', 'es', 'pt-br', 'zh-cn', 'he', 'ar'];
  navigation: NavigationItem[] = [
    // {link: 'about', label: 'anms.menu.about'},
    // {link: 'feature-list', label: 'anms.menu.features'},
    // {link: 'examples', label: 'anms.menu.examples'},
    {link: 'me', label: '我的'},
    {link: 'check-in', label: '打卡'},
    {link: CHECK_IN_TOP_PATH, label: '打卡排行'},
    {link: 'courses', label: '我的课程', roles: ['admin', 'coach']},
    {link: 'admin/courses', label: '课程管理', roles: ['admin']},
    {link: 'admin/users', label: '用户管理', roles: ['admin']}
  ];
  navigationSideMenu = [
    ...this.navigation
    // {link: 'settings', label: 'anms.menu.settings'}
  ];

  isAuthenticated$: Observable<boolean> | undefined;
  stickyHeader$: Observable<boolean> | undefined;
  language$: Observable<string> | undefined;
  theme$: Observable<string> | undefined;
  public user: UserClass | undefined;

  constructor(
    private store: Store<AppState>,
    private storageService: LocalStorageService,
    private authService: AuthService,
    private activeRoute: ActivatedRoute,
    private router: Router
  ) {
    this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        const sp = event.url.split('?');
        console.log('route event', event, sp[0]);
        this.authLayout = sp[0].startsWith('/login');
      }
    });
  }

  private static isIEorEdgeOrSafari() {
    return ['ie', 'edge', 'safari'].includes(browser().name || '');
  }

  ngOnInit(): void {
    this.storageService.testLocalStorage();
    if (AppComponent.isIEorEdgeOrSafari()) {
      this.store.dispatch(
        actionSettingsChangeAnimationsPageDisabled({
          pageAnimationsDisabled: true
        })
      );
    }
    this.isAuthenticated$ = this.authService.isAuthenticated$;
    this.stickyHeader$ = this.store.pipe(select(selectSettingsStickyHeader));
    this.language$ = this.store.pipe(select(selectSettingsLanguage));
    this.theme$ = this.store.pipe(select(selectEffectiveTheme));
    this.authService.user$.subscribe((user) => (this.user = user));
  }

  onLoginClick() {
    this.router.navigate(['/login'], {queryParams: {type: 'login'}});
  }

  onLogoutClick() {
    this.authService.logout();
    location.reload();
  }

  onLanguageSelect(event: MatSelectChange) {
    this.store.dispatch(
      actionSettingsChangeLanguage({language: event.value})
    );
  }

  onProfileClick() {
    this.router.navigate(['/profile', this.authService.getUser().uuid]);
  }

  onSettingClick() {
    this.router.navigate(['/setting']);
  }
}
