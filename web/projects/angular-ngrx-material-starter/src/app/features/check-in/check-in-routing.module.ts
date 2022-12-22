import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { CheckInComponent } from './check-in.component';
import { CheckInShowComponent } from './check-in-show/check-in-show.component';
import { AuthService } from '../../core/auth/auth.service';

const routes: Routes = [
  {
    path: '',
    component: CheckInComponent,
    data: {title: 'anms.menu.check-in'},
    canActivate: [AuthService]
  }, {
    path: ':key',
    component: CheckInShowComponent,
    canActivate: []
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class CheckInRoutingModule {
}
