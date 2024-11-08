import { Component, Input, OnInit } from '@angular/core';
import { Course, CourseService } from '../../../../services/course.service';
import { DatePipe } from '@angular/common';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog } from '@angular/material/dialog';
import { CourseCreateDialogComponent } from '../../courses/course-create-dialog/course-create-dialog.component';

@Component({
  selector: 'app-course-list',
  templateUrl: './course-list.component.html',
  standalone: true,
  imports: [
    DatePipe,
    MatTableModule,
    MatButtonModule
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
    private dialog: MatDialog,
  ) {
  }

  ngOnInit(): void {
    this.load()
  }

  load() {
    this.courseService.getCourseList(this.userId).subscribe(v => {
      this.courses.data = v
    })
  }

  openCreateCourseDialog() {
    this.dialog.open(CourseCreateDialogComponent, {
      minWidth: '400px'
    }).afterClosed().subscribe(v => {
      if (v) {
        this.load()
      }
    })
  }
}
