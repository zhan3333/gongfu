import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  OnInit,
  ViewChild
} from '@angular/core';
import { AdminApiService } from '../../../api/admin/admin-api.service';
import { NotificationService } from '../../../core/notifications/notification.service';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { User } from '../../../api/models/user';
import { Course } from '../../../api/models/course';
import { faEdit, faTrash } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'anms-courses',
  templateUrl: './courses.component.html',
  styleUrls: ['./courses.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
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
    private notification: NotificationService
  ) {}

  ngOnInit(): void {}

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
    this.paginator.pageSize = 10;
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
}
