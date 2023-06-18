import { NgModule } from '@angular/core';
import { PreloadAllModules, RouterModule, Routes } from '@angular/router';
import { AuthService } from './core/auth/auth.service';
import { CheckInTopComponent } from './features/check-in/check-in-top/check-in-top.component';
import { CheckInHistoriesComponent } from './features/check-in/check-in-histories/check-in-histories.component';
import {
  CHECK_IN_HISTORIES_PATH,
  CHECK_IN_TOP_PATH
} from './core/router/route-path';

const routes: Routes = [
  {
    path: 'login',
    loadChildren: () =>
      import('./features/login/login.module').then((m) => m.LoginModule),
    data: {
      authLayout: true
    }
  },
  {
    path: '',
    redirectTo: 'about',
    pathMatch: 'full'
  },
  {
    path: 'about',
    loadChildren: () =>
      import('./features/about/about.module').then((m) => m.AboutModule)
  },
  {
    path: 'me',
    loadComponent: () =>
      import('./features/me/me.component').then((m) => m.MeComponent),
    data: { title: 'anms.menu.me' },
    canActivate: [AuthService]
  },
  {
    path: 'feature-list',
    loadChildren: () =>
      import('./features/feature-list/feature-list.module').then(
        (m) => m.FeatureListModule
      )
  },
  {
    path: 'settings',
    loadChildren: () =>
      import('./features/settings/settings.module').then(
        (m) => m.SettingsModule
      )
  },
  {
    path: 'courses',
    loadChildren: () =>
      import('./features/courses/courses.module').then((m) => m.CoursesModule),
    canActivate: [AuthService]
  },
  {
    path: 'profile/:uuid',
    loadChildren: () =>
      import('./features/profile/profile.module').then((m) => m.ProfileModule)
  },
  {
    path: 'check-in',
    loadChildren: () =>
      import('./features/check-in/check-in.module').then((m) => m.CheckInModule)
  },
  {
    path: CHECK_IN_TOP_PATH,
    component: CheckInTopComponent,
    canActivate: []
  },
  {
    path: CHECK_IN_HISTORIES_PATH,
    component: CheckInHistoriesComponent,
    canActivate: []
  },
  {
    path: 'setting',
    loadChildren: () =>
      import('./features/setting/setting.module').then((m) => m.SettingModule),
    canActivate: [AuthService]
  },
  {
    path: 'admin',
    loadChildren: () =>
      import('./features/admin/admin.module').then((m) => m.AdminModule),
    canActivate: [AuthService]
  },
  {
    path: '**',
    redirectTo: 'me'
  }
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, {
      scrollPositionRestoration: 'enabled',
      preloadingStrategy: PreloadAllModules
    })
  ],
  exports: [RouterModule]
})
export class AppRoutingModule {}
