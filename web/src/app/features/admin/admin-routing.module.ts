import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  {
    path: '',
    children: [
      {
        path: 'users/:id',
        loadComponent: () =>
          import('./user-edit/user-edit.component').then(
            (m) => m.UserEditComponent
          ),
        data: { title: 'user edit' }
      },
      {
        path: 'courses/create',
        loadComponent: () =>
          import('./courses/course-create/course-create.component').then(
            (m) => m.CourseCreateComponent
          ),
        data: { title: 'create course' }
      },
      {
        path: 'courses/:id',
        loadComponent: () =>
          import('./courses/course-edit/course-edit.component').then(
            (m) => m.CourseEditComponent
          ),
        data: { title: 'edit course' }
      },
      {
        path: 'courses',
        loadComponent: () =>
          import('./courses/courses.component').then((m) => m.CoursesComponent),
        data: { title: 'courses page' }
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AdminRoutingModule {}
