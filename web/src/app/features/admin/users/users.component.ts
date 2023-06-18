import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  OnInit,
  ViewChild
} from '@angular/core';
import { MatTableDataSource } from '@angular/material/table';
import { User } from '../../../api/models/user';
import { MatPaginator } from '@angular/material/paginator';
import { AdminApiService } from '../../../api/admin/admin-api.service';
import { faEdit } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'anms-users',
  templateUrl: './users.component.html',
  changeDetection: ChangeDetectionStrategy.Default
})
export class UsersComponent implements OnInit, AfterViewInit {
  @ViewChild(MatPaginator) paginator: MatPaginator | undefined;
  public faEdit = faEdit;

  public displayedColumns: string[] = [
    'actions',
    'avatar',
    'nickname',
    'phone',
    'roles'
  ];
  public dataSource = new MatTableDataSource<User>([]);

  constructor(private adminApi: AdminApiService) {}

  ngOnInit(): void {}

  refreshTable() {
    if (this.paginator === undefined) {
      return;
    }
    this.adminApi
      .getUsers({
        desc: true,
        keyword: '',
        limit: this.paginator.pageSize,
        page: this.paginator.pageIndex
      })
      .subscribe((data) => {
        this.dataSource.data = data.users;
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
}
