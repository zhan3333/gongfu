import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { CoursesComponent } from './courses.component';
import { CourseCreateComponent } from '../admin/courses/course-create/course-create.component';
import { CourseEditComponent } from './course-edit/course-edit.component';

const routes: Routes = [
  {
    path: '',
    component: CoursesComponent,
    data: {title: 'anms.menu.courses'}
  }
  ,
  {
    path: ':id',
    component: CourseEditComponent,
    data: {title: 'anms.menu.course-edit'}
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class CoursesRoutingModule {
}
