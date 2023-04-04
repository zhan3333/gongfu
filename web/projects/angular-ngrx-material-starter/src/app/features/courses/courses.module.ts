import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CoursesComponent } from './courses.component';
import { SharedModule } from '../../shared/shared.module';
import { CoursesRoutingModule } from './courses-routing.module';
import { CourseEditComponent } from './course-edit/course-edit.component';


@NgModule({
  declarations: [
    CoursesComponent,
    CourseEditComponent
  ],
  imports: [
    CommonModule, SharedModule, CoursesRoutingModule
  ]
})
export class CoursesModule {
}
