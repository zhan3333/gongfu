import { NgModule } from '@angular/core';
import { SharedModule } from '../../shared/shared.module';
import { UsersComponent } from './users/users.component';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyPaginatorModule as MatPaginatorModule } from '@angular/material/legacy-paginator';
import { AdminRoutingModule } from './admin-routing.module';
import { UserEditComponent } from './user-edit/user-edit.component';
import { CoursesComponent } from './courses/courses.component';
import { CourseCreateComponent } from './courses/course-create/course-create.component';
import { CourseEditComponent } from './courses/course-edit/course-edit.component';

@NgModule({
  declarations: [
    UsersComponent,
    UserEditComponent,
    CoursesComponent,
    CourseCreateComponent,
    CourseEditComponent
  ],
  imports: [
    SharedModule,
    MatTableModule,
    MatPaginatorModule,
    AdminRoutingModule
  ]
})
export class AdminModule {}
