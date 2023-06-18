import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { UsersComponent } from './users/users.component';
import { UserEditComponent } from './user-edit/user-edit.component';
import { CourseCreateComponent } from './courses/course-create/course-create.component';
import { CoursesComponent } from './courses/courses.component';
import { CourseEditComponent } from './courses/course-edit/course-edit.component';

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
      {
        path: 'courses/create',
        component: CourseCreateComponent,
        data: {title: 'create course'}
      },
      {
        path: 'courses/:id',
        component: CourseEditComponent,
        data: {title: 'edit course'}
      },
      {
        path: 'courses',
        component: CoursesComponent,
        data: {title: 'courses page'}
      }
    ]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AdminRoutingModule {
}
