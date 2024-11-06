import { Component, Input, OnInit } from '@angular/core';
import { Course, CourseService } from '../../../../services/course.service';
import { DatePipe } from '@angular/common';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';

@Component({
  selector: 'app-course-list',
  templateUrl: './course-list.component.html',
  standalone: true,
  imports: [
    DatePipe,
    MatTableModule
  ]
})
export class CourseListComponent implements OnInit {
  @Input() userId!: number;
  public courses = new MatTableDataSource<Course>([]);
  public displayedColumns: string[] = [
    "id",
    "schoolName",
    "startDate",
    "managerName",
    "coachName",
    "assistantNames",
    "checkInAt",
    "checkOutAt",
    "summary",
    "content",
  ]

  constructor(
    private courseService: CourseService,
  ) {
  }

  ngOnInit(): void {
    this.courseService.getCourseList(this.userId).subscribe(v => {
      this.courses.data = v
    })
  }
}
