import { NgModule } from '@angular/core';
import { Routes, RouterModule, PreloadAllModules } from '@angular/router';
import { AuthService } from './core/auth/auth.service';
import { CheckInShowComponent } from './features/check-in/check-in-show/check-in-show.component';
import { CheckInTopComponent } from './features/check-in/check-in-top/check-in-top.component';
import { CheckInCountComponent } from './features/check-in/check-in-count/check-in-count.component';
import { CheckInHistoriesComponent } from './features/check-in/check-in-histories/check-in-histories.component';
import { CheckInContinuousComponent } from './features/check-in/check-in-continuous/check-in-continuous.component';
import {
  CHECK_IN_CONTINUOUS_TOP_PATH,
  CHECK_IN_COUNT_PATH,
  CHECK_IN_HISTORIES_PATH,
  CHECK_IN_TOP_PATH
} from './core/router/route-path';

const routes: Routes = [
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
    path: 'me',
    loadChildren: () =>
      import('./features/me/me.module').then(
        (m) => m.MeModule
      )
  },
  {
    path: 'check-in',
    loadChildren: () =>
      import('./features/check-in/check-in.module').then(
        (m) => m.CheckInModule
      ),
    canActivate: [AuthService]
  },
  {
    path: CHECK_IN_TOP_PATH,
    component: CheckInTopComponent,
    canActivate: []
  },
  {
    path: CHECK_IN_COUNT_PATH,
    component: CheckInCountComponent,
    canActivate: []
  },
  {
    path: CHECK_IN_CONTINUOUS_TOP_PATH,
    component: CheckInContinuousComponent,
    canActivate: []
  },
  {
    path: CHECK_IN_HISTORIES_PATH,
    component: CheckInHistoriesComponent,
    canActivate: []
  },
  {
    path: 'check-in/:key',
    component: CheckInShowComponent,
    canActivate: []
  },
  {
    path: 'examples',
    loadChildren: () =>
      import('./features/examples/examples.module').then(
        (m) => m.ExamplesModule
      )
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
      preloadingStrategy: PreloadAllModules,
      relativeLinkResolution: 'legacy'
    })
  ],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
