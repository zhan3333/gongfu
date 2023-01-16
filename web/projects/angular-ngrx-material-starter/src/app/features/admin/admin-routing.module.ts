import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { UsersComponent } from './users/users.component';
import { UserEditComponent } from './user-edit/user-edit.component';

const routes: Routes = [
  {
    path: '',
    children: [
      {
        path: 'users',
        component: UsersComponent,
        data: {title: 'users'},
      },
      {
        path: 'users/:id',
        component: UserEditComponent,
        data: {title: 'user edit'},
      },
    ]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AdminRoutingModule {
}
