import browser from 'browser-detect';
import { Component, OnInit } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { Store, select } from '@ngrx/store';
import { Observable } from 'rxjs';

import { environment as env } from '../../environments/environment';

import {
  authLogin,
  authLogout,
  routeAnimations,
  LocalStorageService,
  selectIsAuthenticated,
  selectSettingsStickyHeader,
  selectSettingsLanguage,
  selectEffectiveTheme,
  AppState
} from '../core/core.module';
import {
  actionSettingsChangeAnimationsPageDisabled,
  actionSettingsChangeLanguage
} from '../core/settings/settings.actions';
import { AuthService } from '../core/auth/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'anms-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  animations: [routeAnimations]
})
export class AppComponent implements OnInit {
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
    {link: 'check-in/top', label: '打卡排行'}
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
    private router: Router
  ) {
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
    this.isAuthenticated$.subscribe(b => console.log('isAuthenticated$', b))

    this.stickyHeader$ = this.store.pipe(select(selectSettingsStickyHeader));
    this.language$ = this.store.pipe(select(selectSettingsLanguage));
    this.theme$ = this.store.pipe(select(selectEffectiveTheme));
  }

  onLoginClick() {
    window.location.href = '/login'
  }

  onLogoutClick() {
    this.authService.logout()
  }

  onLanguageSelect(event: MatSelectChange) {
    this.store.dispatch(
      actionSettingsChangeLanguage({language: event.value})
    );
  }
}
