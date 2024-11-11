import { AfterViewInit, ChangeDetectionStrategy, Component, OnInit, ViewChild } from '@angular/core';
import { AdminApiService } from '../../../api/admin/admin-api.service';
import { NotificationService } from '../../../core/notifications/notification.service';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { Course } from '../../../api/models/course';
import { faEdit, faTrash } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { MatIconModule } from '@angular/material/icon';
import { RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog } from '@angular/material/dialog';
import { CourseCreateDialogComponent } from './course-create-dialog/course-create-dialog.component';

@Component({
  selector: 'anms-courses',
  templateUrl: './courses.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [MatButtonModule, RouterLink, MatTableModule, MatIconModule, FontAwesomeModule, MatPaginatorModule]
})
export class CoursesComponent implements OnInit, AfterViewInit {
  @ViewChild(MatPaginator) paginator: MatPaginator | undefined;
  public faTrash = faTrash;
  public faEdit = faEdit;

  public displayedColumns: string[] = [
    'actions',
    'id',
    'school',
    'coach',
    'startDateTime'
  ];
  public dataSource = new MatTableDataSource<Course>([]);

  constructor(
    private adminApi: AdminApiService,
    private notification: NotificationService,
    private dialog: MatDialog,
  ) {
  }

  ngOnInit(): void {
  }

  refreshTable() {
    if (this.paginator === undefined) {
      return;
    }
    this.adminApi
      .getCoursesPage({
        desc: true,
        keyword: '',
        limit: this.paginator.pageSize,
        page: this.paginator.pageIndex
      })
      .subscribe((data) => {
        this.dataSource.data = data.items;
        if (this.paginator === undefined) {
          return;
        }
        this.paginator.pageSize = data.limit;
        this.paginator.pageIndex = data.page;
        this.paginator.length = data.count;
      });
  }

  ngAfterViewInit(): void {
    if (this.paginator === undefined) {
      return;
    }
    // this.paginator.pageSize = 10;
    this.paginator.page.subscribe((v) => {
      this.refreshTable();
    });
    // this.dataSource.paginator = this.paginator
    this.refreshTable();
  }

  // 删除课程
  onClickDeleteCourse(id: number) {
    this.adminApi.deleteCourse(id).subscribe(() => {
      this.refreshTable();
      this.notification.success('删除成功');
    });
  }

  openCreateCourseDialog() {
    this.dialog.open(CourseCreateDialogComponent, {
      minWidth: '400px'
    }).afterClosed().subscribe(v => {
      if (v) {
        this.refreshTable();
      }
    })
  }
}
