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
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';
import { CHECK_IN_TOP_PATH } from '../core/router/route-path';

@Component({
  selector: 'anms-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  animations: [routeAnimations],
  changeDetection: ChangeDetectionStrategy.Default
})
export class AppComponent implements OnInit {
  authLayout = false;
  displayFooter = false;
  isProd = env.production;
  envName = env.envName;
  version = env.versions.app;
  year = new Date().getFullYear();
  logo = 'assets/logo.png';
  languages = ['en', 'de', 'sk', 'fr', 'es', 'pt-br', 'zh-cn', 'he', 'ar'];
  navigation = [
    // {link: 'about', label: 'anms.menu.about'},
    // {link: 'feature-list', label: 'anms.menu.features'},
    // {link: 'examples', label: 'anms.menu.examples'},
    {link: 'me', label: '我的'},
    {link: 'check-in', label: '打卡'},
    {link: CHECK_IN_TOP_PATH, label: '打卡排行'},
  ];
  navigationSideMenu = [
    ...this.navigation,
    // {link: 'settings', label: 'anms.menu.settings'}
  ];

  isAuthenticated$: Observable<boolean> | undefined;
  stickyHeader$: Observable<boolean> | undefined;
  language$: Observable<string> | undefined;
  theme$: Observable<string> | undefined;

  constructor(
    private store: Store<AppState>,
    private storageService: LocalStorageService,
    private authService: AuthService,
    private activeRoute: ActivatedRoute,
    private router: Router,
  ) {
    this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        const sp = event.url.split('?')
        console.log('route event', event, sp[0])
        this.authLayout = sp[0].startsWith('/login');
      }
    })
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
    if (this.authService.isAuthenticated() && !this.authService.getUser().id) {
      this.authService.logout()
      return
    }
  }

  onLoginClick() {
    this.router.navigate(['/login'], {queryParams: {'type': 'login'}})
  }

  onLogoutClick() {
    this.authService.logout()
    location.reload()
  }

  onLanguageSelect(event: MatSelectChange) {
    this.store.dispatch(
      actionSettingsChangeLanguage({language: event.value})
    );
  }

  onProfileClick() {
    this.router.navigate(['/profile', this.authService.getUser().uuid])
  }

  onSettingClick() {
    this.router.navigate(['/setting'])
  }
}
