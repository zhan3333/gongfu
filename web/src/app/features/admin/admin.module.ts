import { NgModule } from '@angular/core';
import { SharedModule } from '../../shared/shared.module';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule } from '@angular/material/paginator';
import { AdminRoutingModule } from './admin-routing.module';
import { UserEditComponent } from './user-edit/user-edit.component';
import { CoursesComponent } from './courses/courses.component';
import { CourseCreateComponent } from './courses/course-create/course-create.component';
import { CourseEditComponent } from './courses/course-edit/course-edit.component';

@NgModule({
  imports: [
    SharedModule,
    MatTableModule,
    MatPaginatorModule,
    AdminRoutingModule,
    UserEditComponent,
    CoursesComponent,
    CourseCreateComponent,
    CourseEditComponent
  ],
  declarations: []
})
export class AdminModule {
}
